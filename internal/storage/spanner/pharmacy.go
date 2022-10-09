package spanner

import (
	spannerSyntax "cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/justdomepaul/toolbox/database/spanner"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/spannertool"
	"github.com/justdomepaul/toolbox/stringtool"
	"github.com/justdomepaul/toolbox/timestamp"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"math"
	"phantom_mask/internal/entity"
	"phantom_mask/internal/storage"
	"time"
)

var (
	pharmacyTable = "Pharmacy"
)

var pharmacyClauseFn = map[storage.PharmacyEnumType]func(source storage.PharmacyListCondition, condition *string, args map[string]interface{}) error{
	storage.PharmacyProductPriceRange: withPharmacyProductPriceRange,
}

func withPharmacyProductPriceRange(source storage.PharmacyListCondition, condition *string, args map[string]interface{}) error {
	*condition = stringtool.StringJoin(*condition, ` AND @Min <= price AND price <= @Max`)
	args["Min"] = source.Min
	args["Max"] = source.Max
	return nil
}

func toPharmacyClauses(source storage.PharmacyListCondition) (conditionSyntax string, args map[string]interface{}, err error) {
	args = map[string]interface{}{}
	for _, op := range source.Fields {
		if err := pharmacyClauseFn[op](source, &conditionSyntax, args); err != nil {
			return conditionSyntax, args, err
		}
	}
	return conditionSyntax, args, err
}

// NewPharmacy method
func NewPharmacy(logger *zap.Logger, session spanner.ISession) *Pharmacy {
	return &Pharmacy{
		logger:  logger,
		session: session,
	}
}

type Pharmacy struct {
	logger  *zap.Logger
	session spanner.ISession
}

func (st Pharmacy) Create(ctx context.Context, input entity.Pharmacy) error {
	if err := validator.New().Struct(&input); err != nil {
		return fmt.Errorf("%w: %s", errorhandler.ErrInvalidArguments, err.Error())
	}
	_, err := st.session.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spannerSyntax.ReadWriteTransaction) error {
		columns, placeholder, params := spannertool.FetchSpannerTagValue(input, false, DBCreatedTime)
		stmt := spannerSyntax.Statement{
			SQL:    fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, pharmacyTable, columns, placeholder),
			Params: params,
		}
		_, err := txn.Update(ctx, stmt)
		return err
	})
	if spannerSyntax.ErrCode(err) == codes.AlreadyExists {
		return fmt.Errorf("%w: %s", errorhandler.ErrAlreadyExists, err.Error())
	}
	return err
}

func (st Pharmacy) ListPharmacyMixProduct(ctx context.Context, row, page uint64, name string, orderEnum storage.OrderListEnum) (*entity.PharmacyProductList, error) {
	if err := spannertool.ValidListArgument(row, page); err != nil {
		return nil, err
	}

	args := map[string]interface{}{}

	args["Row"] = int64(row)
	args["Page"] = int64(page)
	args["Offset"] = int64((page - 1) * row)

	args["Name"] = fmt.Sprintf(`\Q%s\E`, name)

	stmt := spannerSyntax.Statement{
		SQL: fmt.Sprintf(
			`
WITH Data AS (
    SELECT Ph.UID AS UID, Ph.Name AS PharmacyName, CashBalance, ProductID, P.Name AS ProductName, Price 
	FROM %s AS Ph JOIN %s AS P on Ph.UID = P.UID 
	WHERE REGEXP_CONTAINS(Ph.Name, @Name) OR REGEXP_CONTAINS(P.Name, @Name)
)
SELECT 
	(SELECT COUNT(*) FROM Data) AS Count, 
	@Row AS Row, 
	@Page AS Page, 
	(SELECT ARRAY(
		SELECT STRUCT(UID, ProductID, PharmacyName, CashBalance, ProductName, Price) 
		FROM Data%s LIMIT @Row OFFSET @Offset
	)) AS PharmacyProducts
`, pharmacyTable, productTable, withTimeOrder(orderEnum),
		),
		Params: args,
	}
	iter := st.session.Single().Query(ctx, stmt)
	defer iter.Stop()

	resp := &entity.PharmacyProductList{}
	if err := spannertool.GetIteratorFirstRow(iter, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (st Pharmacy) ListSpecifyTime(ctx context.Context, row, page uint64, specifyTimestamp int64, orderEnum storage.OrderListEnum) (*entity.PharmacySpecifyTimestampList, error) {
	if err := spannertool.ValidListArgument(row, page); err != nil {
		return nil, err
	}

	args := map[string]interface{}{}

	args["Row"] = int64(row)
	args["Page"] = int64(page)
	args["Offset"] = int64((page - 1) * row)

	specify := timestamp.GetUTC8Time(specifyTimestamp)
	specifyHour, err := time.ParseDuration(fmt.Sprintf("%dh%dm", specify.Hour(), specify.Minute()))
	if err != nil {
		return nil, err
	}

	args["SpecifyDay"] = int64(specify.Weekday())
	args["SpecifyTime"] = math.Round(specifyHour.Hours()*100) / 100

	stmt := spannerSyntax.Statement{
		SQL: fmt.Sprintf(
			`
WITH SpecifyTimeData AS (
    SELECT
        P.UID AS UID, Name, CashBalance, CreatedTime, Day, OpenHour, 
		(CASE WHEN CloseHour > 24 THEN CloseHour-24 ELSE CloseHour END) AS CloseHour,
        (CASE WHEN @SpecifyTime < OpenHour
        THEN CASE WHEN OpenHour <= @SpecifyTime + 24 AND @SpecifyTime + 24 < CloseHour THEN true ELSE false END
        ELSE CASE WHEN OpenHour <= @SpecifyTime AND @SpecifyTime < CloseHour THEN true ELSE false END
        END) AS inRange
    FROM %s AS P JOIN %s AS PI on P.UID = PI.UID WHERE Day = @SpecifyDay
)
SELECT 
	(SELECT COUNT(*) FROM SpecifyTimeData WHERE inRange = true) AS Count, 
	@Row AS Row, 
	@Page AS Page, 
	(SELECT ARRAY(
		SELECT STRUCT(UID, Name, CashBalance, CreatedTime, Day, OpenHour, CloseHour) 
		FROM specifyTimeData WHERE inRange = true%s LIMIT @Row OFFSET @Offset
	)) AS Pharmacies
`, pharmacyTable, pharmacyInfoTable, withTimeOrder(orderEnum),
		),
		Params: args,
	}
	iter := st.session.Single().Query(ctx, stmt)
	defer iter.Stop()

	resp := &entity.PharmacySpecifyTimestampList{}
	if err := spannertool.GetIteratorFirstRow(iter, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (st Pharmacy) ListByProductPriceRange(ctx context.Context, row, page uint64, orderEnum storage.OrderListEnum, condition storage.PharmacyListCondition) (*entity.PharmacyList, error) {
	if err := spannertool.ValidListArgument(row, page); err != nil {
		return nil, err
	}

	conditionSyntax, args, err := toPharmacyClauses(condition)
	if err != nil {
		return nil, err
	}

	args["Row"] = int64(row)
	args["Offset"] = int64((page - 1) * row)
	args["Page"] = int64(page)

	stmt := spannerSyntax.Statement{
		SQL: fmt.Sprintf(
			`
WITH Data AS (
    SELECT DISTINCT Ph.* FROM %s AS Ph JOIN %s P on Ph.UID = P.UID WHERE Ph.UID IS NOT NULL%s
)
SELECT 
	(SELECT COUNT(*) FROM Data) AS Count, 
	@Row AS Row, 
	@Page AS Page, 
	(SELECT ARRAY(
		SELECT STRUCT(UID, Name, CashBalance, CreatedTime) 
		FROM Data%s LIMIT @Row OFFSET @Offset
	)) AS Pharmacies
`, pharmacyTable, productTable, conditionSyntax, withTimeOrder(orderEnum),
		),
		Params: args,
	}
	iter := st.session.Single().Query(ctx, stmt)
	defer iter.Stop()

	resp := &entity.PharmacyList{}
	if err := spannertool.GetIteratorFirstRow(iter, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
