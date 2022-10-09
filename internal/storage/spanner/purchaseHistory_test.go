package spanner

import (
	"context"
	"github.com/google/uuid"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"phantom_mask/internal/entity"
	"testing"
	"time"
)

type PurchaseHistorySuite struct {
	suite.Suite
	ctx         context.Context
	logger      *zap.Logger
	client      *PurchaseHistory
	userUID     []byte
	pharmacyUID []byte
	productUID  []byte
}

func (suite *PurchaseHistorySuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, err := zap.NewDevelopment()
	suite.NoError(err)
	suite.logger = logger
	suite.client = NewPurchaseHistory(suite.logger, session)

	userUID, err := uuid.NewUUID()
	suite.NoError(err)
	pharmacyUID, err := uuid.NewUUID()
	suite.NoError(err)
	productUID, err := uuid.NewUUID()
	suite.NoError(err)
	userClient := NewUser(suite.logger, session)
	pharmacyClient := NewPharmacy(suite.logger, session)
	productClient := NewProduct(suite.logger, session)
	suite.NoError(userClient.Create(suite.ctx, entity.User{
		UID:         userUID[:],
		Name:        "Yvonne Guerrero",
		CashBalance: 10.5,
	}))
	suite.userUID = userUID[:]
	suite.NoError(pharmacyClient.Create(suite.ctx, entity.Pharmacy{
		UID:         pharmacyUID[:],
		Name:        "Carepoint",
		CashBalance: 10.5,
	}))
	suite.pharmacyUID = pharmacyUID[:]
	suite.NoError(productClient.Create(suite.ctx, entity.Product{
		UID:       pharmacyUID[:],
		ProductID: productUID[:],
		Name:      "True Barrier (green) (3 per pack)",
		Price:     10.5,
	}))
	suite.productUID = productUID[:]
}

func (suite *PurchaseHistorySuite) TestCreateMethod() {
	type want struct {
		Error error
	}

	specifyTime, err := time.Parse("2006-01-02 15:04:05", "2021-01-04 15:18:51")
	suite.NoError(err)

	testCases := []struct {
		Label  string
		Entity entity.PurchaseHistory
		Want   want
	}{
		{
			Label: "CreatePharmacyShouldSuccess",
			Entity: entity.PurchaseHistory{
				UID:               suite.userUID,
				PharmacyUID:       suite.pharmacyUID,
				ProductID:         suite.productUID,
				TransactionAmount: 12.35,
				TransactionDate:   specifyTime,
			},
			Want: want{},
		},
		{
			Label: "CreateDuplicatePharmacyShouldFail",
			Entity: entity.PurchaseHistory{
				UID:               suite.userUID,
				PharmacyUID:       suite.pharmacyUID,
				ProductID:         suite.productUID,
				TransactionAmount: 12.35,
				TransactionDate:   specifyTime,
			},
			Want: want{
				Error: errorhandler.ErrAlreadyExists,
			},
		},
	}

	for _, tc := range testCases {
		if tc.Want.Error != nil {
			suite.ErrorIs(suite.client.Create(suite.ctx, tc.Entity), tc.Want.Error)
		} else {
			suite.NoError(suite.client.Create(suite.ctx, tc.Entity))
		}
	}
}

func TestPurchaseHistorySuite(t *testing.T) {
	suite.Run(t, new(PurchaseHistorySuite))
}
