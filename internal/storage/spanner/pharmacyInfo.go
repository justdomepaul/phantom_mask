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
)

var (
	pharmacyInfoTable = "PharmacyInfo"
)

// NewPharmacyInfo method
func NewPharmacyInfo(logger *zap.Logger, session spanner.ISession) *PharmacyInfo {
	return &PharmacyInfo{
		logger:  logger,
		session: session,
	}
}

type PharmacyInfo struct {
	logger  *zap.Logger
	session spanner.ISession
}

func (st PharmacyInfo) Create(ctx context.Context, input entity.PharmacyInfo) error {
	if err := validator.New().Struct(&input); err != nil {
		return fmt.Errorf("%w: %s", errorhandler.ErrInvalidArguments, err.Error())
	}
	_, err := st.session.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spannerSyntax.ReadWriteTransaction) error {
		columns, placeholder, params := spannertool.FetchSpannerTagValue(input, false, DBCreatedTime)
		stmt := spannerSyntax.Statement{
			SQL:    fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, pharmacyInfoTable, columns, placeholder),
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
