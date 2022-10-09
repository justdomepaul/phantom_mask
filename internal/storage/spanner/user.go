package spanner

import (
	spannerSyntax "cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/justdomepaul/toolbox/database/spanner"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/spannertool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"phantom_mask/internal/entity"
	"time"
)

var (
	userTable = "User"
)

// NewUser method
func NewUser(logger *zap.Logger, session spanner.ISession) *User {
	return &User{
		logger:  logger,
		session: session,
	}
}

type User struct {
	logger  *zap.Logger
	session spanner.ISession
}

func (st User) Create(ctx context.Context, input entity.User) error {
	if err := validator.New().Struct(&input); err != nil {
		return fmt.Errorf("%w: %s", errorhandler.ErrInvalidArguments, err.Error())
	}
	_, err := st.session.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spannerSyntax.ReadWriteTransaction) error {
		columns, placeholder, params := spannertool.FetchSpannerTagValue(input, false, DBCreatedTime)
		stmt := spannerSyntax.Statement{
			SQL:    fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, userTable, columns, placeholder),
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

func (st User) ListTopTransactionAmount(ctx context.Context, topNumber, startTime, endTime int64) (*entity.TopTransactionAmountList, error) {
	stmt := spannerSyntax.Statement{
		SQL: fmt.Sprintf(
			`
WITH Data AS (
    SELECT U.UID, Name, SUM(TransactionAmount) AS TransactionAmount
    FROM %s AS U JOIN %s PH on U.UID = PH.UID
    WHERE @StartTime <= TransactionDate AND TransactionDate <= @EndTime
    GROUP BY U.UID, Name
)
SELECT
    (SELECT ARRAY(
        SELECT STRUCT(UID, Name, TransactionAmount) FROM Data ORDER BY TransactionAmount DESC LIMIT @TopNumber
    )) AS TopTransactionAmountUsers
`, userTable, purchaseHistoryTable),
		Params: map[string]interface{}{
			"StartTime": time.UnixMilli(startTime),
			"EndTime":   time.UnixMilli(endTime),
			"TopNumber": topNumber,
		},
	}
	iter := st.session.Single().Query(ctx, stmt)
	defer iter.Stop()

	resp := &entity.TopTransactionAmountList{}
	if err := spannertool.GetIteratorFirstRow(iter, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (st User) GetTransactionTotal(ctx context.Context, startTime, endTime int64) (*entity.TransactionTotal, error) {
	stmt := spannerSyntax.Statement{
		SQL: fmt.Sprintf(
			`
WITH Data AS (
    SELECT ProductID, SUM(TransactionAmount) AS TransactionAmount FROM %s 
	WHERE @StartTime <= TransactionDate AND TransactionDate <= @EndTime GROUP BY ProductID
)
SELECT
    (SELECT COUNT(ProductID) FROM Data) AS Total,
    (SELECT SUM(TransactionAmount) FROM Data) AS TransactionAmount
`, purchaseHistoryTable),
		Params: map[string]interface{}{
			"StartTime": time.UnixMilli(startTime),
			"EndTime":   time.UnixMilli(endTime),
		},
	}
	iter := st.session.Single().Query(ctx, stmt)
	defer iter.Stop()

	resp := &entity.TransactionTotal{}
	if err := spannertool.GetIteratorFirstRow(iter, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
