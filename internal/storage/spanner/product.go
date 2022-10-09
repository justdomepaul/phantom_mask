package spanner

import (
	spannerSyntax "cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/justdomepaul/toolbox/database/spanner"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/spannertool"
	"github.com/justdomepaul/toolbox/stringtool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"phantom_mask/internal/entity"
	"phantom_mask/internal/storage"
	"time"
)

var (
	productTable = "Product"
)

var productClauseFn = map[storage.ProductEnumType]func(source storage.ProductListCondition, condition *string, args map[string]interface{}) error{
	storage.ProductSpecifyPharmacy: withProductSpecifyPharmacy,
}

func withProductSpecifyPharmacy(source storage.ProductListCondition, condition *string, args map[string]interface{}) error {
	if err := validator.New().Var(source.PharmacyID, `required`); err != nil {
		return err
	}
	*condition = stringtool.StringJoin(*condition, ` AND UID = @PharmacyID`)
	args["PharmacyID"] = source.PharmacyID
	return nil
}

func toProductClauses(source storage.ProductListCondition) (conditionSyntax string, args map[string]interface{}, err error) {
	args = map[string]interface{}{}
	for _, op := range source.Fields {
		if err := productClauseFn[op](source, &conditionSyntax, args); err != nil {
			return conditionSyntax, args, err
		}
	}
	return conditionSyntax, args, err
}

// NewProduct method
func NewProduct(logger *zap.Logger, session spanner.ISession) *Product {
	return &Product{
		logger:  logger,
		session: session,
	}
}

type Product struct {
	logger  *zap.Logger
	session spanner.ISession
}

func (st Product) Create(ctx context.Context, input entity.Product) error {
	if err := validator.New().Struct(&input); err != nil {
		return fmt.Errorf("%w: %s", errorhandler.ErrInvalidArguments, err.Error())
	}
	_, err := st.session.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spannerSyntax.ReadWriteTransaction) error {
		columns, placeholder, params := spannertool.FetchSpannerTagValue(input, false, DBCreatedTime)
		stmt := spannerSyntax.Statement{
			SQL:    fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, productTable, columns, placeholder),
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

func (st Product) Purchase(ctx context.Context, userID, pharmacyID, productID []byte, quantity int) error {
	input := struct {
		UserID     []byte `json:"user_id,omitempty" validate:"required"`
		PharmacyID []byte `json:"pharmacy_id,omitempty" validate:"required"`
		ProductID  []byte `json:"product_id,omitempty" validate:"required"`
		Quantity   int    `json:"quantity,omitempty" validate:"required,min=1"`
	}{
		UserID:     userID,
		PharmacyID: pharmacyID,
		ProductID:  productID,
		Quantity:   quantity,
	}
	if err := validator.New().Struct(&input); err != nil {
		return fmt.Errorf("%w: %s", errorhandler.ErrInvalidArguments, err.Error())
	}
	_, err := st.session.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spannerSyntax.ReadWriteTransaction) error {
		getUserBalance := func(key spannerSyntax.Key) (float64, error) {
			row, err := txn.ReadRow(ctx, userTable, key, []string{"CashBalance"})
			if err != nil {
				return 0, err
			}
			var cashBalance float64
			if err := row.Column(0, &cashBalance); err != nil {
				return 0, err
			}
			return cashBalance, nil
		}
		getPharmacyBalance := func(key spannerSyntax.Key) (float64, error) {
			row, err := txn.ReadRow(ctx, pharmacyTable, key, []string{"CashBalance"})
			if err != nil {
				return 0, err
			}
			var cashBalance float64
			if err := row.Column(0, &cashBalance); err != nil {
				return 0, err
			}
			return cashBalance, nil
		}
		getProductPrice := func(key spannerSyntax.Key) (float64, error) {
			row, err := txn.ReadRow(ctx, productTable, key, []string{"Price"})
			if err != nil {
				return 0, err
			}
			var price float64
			if err := row.Column(0, &price); err != nil {
				return 0, err
			}
			return price, nil
		}
		var mut []*spannerSyntax.Mutation
		var (
			userColumns            = []string{"UID", "CashBalance"}
			pharmacyColumns        = []string{"UID", "CashBalance"}
			purchaseHistoryColumns = []string{"UID", "PharmacyUID", "ProductID", "TransactionAmount", "TransactionDate"}
		)
		userCashBalance, err := getUserBalance(spannerSyntax.Key{userID})
		if err != nil {
			return err
		}
		pharmacyBalance, err := getPharmacyBalance(spannerSyntax.Key{pharmacyID})
		if err != nil {
			return err
		}
		productPrice, err := getProductPrice(spannerSyntax.Key{pharmacyID, productID})
		if err != nil {
			return err
		}
		upgradeUserCashBalance := userCashBalance - (productPrice * float64(quantity))
		if upgradeUserCashBalance < 0 {
			return errors.New("user CashBalance not enough")
		}
		mut = append(mut, spannerSyntax.Update(userTable, userColumns, []interface{}{userID, upgradeUserCashBalance}))
		mut = append(mut, spannerSyntax.Update(
			pharmacyTable, pharmacyColumns,
			[]interface{}{pharmacyID, pharmacyBalance + (productPrice * float64(quantity))}))
		mut = append(mut, spannerSyntax.Insert(
			purchaseHistoryTable, purchaseHistoryColumns,
			[]interface{}{userID, pharmacyID, productID, productPrice * float64(quantity), time.Now().UTC()}))

		return txn.BufferWrite(mut)
	})
	if spannerSyntax.ErrCode(err) == codes.AlreadyExists {
		return fmt.Errorf("%w: %s", errorhandler.ErrAlreadyExists, err.Error())
	}
	return err
}

func (st Product) List(ctx context.Context, row, page uint64, orderEnum storage.OrderListEnum, condition storage.ProductListCondition) (*entity.ProductList, error) {
	if err := spannertool.ValidListArgument(row, page); err != nil {
		return nil, err
	}

	conditionSyntax, args, err := toProductClauses(condition)
	if err != nil {
		return nil, err
	}

	args["Row"] = int64(row)
	args["Offset"] = int64((page - 1) * row)
	args["Page"] = int64(page)

	stmt := spannerSyntax.Statement{
		SQL: fmt.Sprintf(
			`
SELECT 
	(SELECT COUNT(*) FROM %s WHERE UID IS NOT NULL%s) AS Count, 
	@Row AS Row, 
	@Page AS Page, 
	(SELECT ARRAY(
		SELECT STRUCT(UID, ProductID, Name, Price, CreatedTime) 
		FROM %s WHERE UID IS NOT NULL%s%s LIMIT @Row OFFSET @Offset
	)) AS Products
`, productTable, conditionSyntax, productTable, conditionSyntax, withTimeOrder(orderEnum),
		),
		Params: args,
	}
	iter := st.session.Single().Query(ctx, stmt)
	defer iter.Stop()

	resp := &entity.ProductList{}
	if err := spannertool.GetIteratorFirstRow(iter, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
