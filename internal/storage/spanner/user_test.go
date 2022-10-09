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

type UserSuite struct {
	suite.Suite
	ctx                   context.Context
	logger                *zap.Logger
	client                *User
	purchaseHistoryClient *PurchaseHistory
	pharmacyClient        *Pharmacy
	productClient         *Product
}

func (suite *UserSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, err := zap.NewDevelopment()
	suite.NoError(err)
	suite.logger = logger
	suite.client = NewUser(suite.logger, session)
	suite.purchaseHistoryClient = NewPurchaseHistory(suite.logger, session)
	suite.pharmacyClient = NewPharmacy(suite.logger, session)
	suite.productClient = NewProduct(suite.logger, session)
}

func (suite *UserSuite) TestCreateMethod() {
	type want struct {
		Error error
	}

	uid, err := uuid.NewUUID()
	suite.NoError(err)

	testCases := []struct {
		Label  string
		Entity entity.User
		Want   want
	}{
		{
			Label: "CreatePharmacyShouldSuccess",
			Entity: entity.User{
				UID:         uid[:],
				Name:        "Yvonne Guerrero",
				CashBalance: 10.5,
			},
			Want: want{},
		},
		{
			Label: "CreateDuplicatePharmacyShouldFail",
			Entity: entity.User{
				UID:         uid[:],
				Name:        "Yvonne Guerrero",
				CashBalance: 10.5,
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

func (suite *UserSuite) TestListTopTransactionAmountMethod() {
	type want struct {
		Error error
	}

	uid, err := uuid.NewUUID()
	suite.NoError(err)
	pharmacyUID, err := uuid.NewUUID()
	suite.NoError(err)
	productUID, err := uuid.NewUUID()
	suite.NoError(err)

	startTime, err := time.Parse("2006-01-02", "2022-10-01")
	suite.NoError(err)
	endTime, err := time.Parse("2006-01-02", "2022-10-10")
	suite.NoError(err)

	suite.NoError(suite.pharmacyClient.Create(suite.ctx, entity.Pharmacy{
		UID:         pharmacyUID[:],
		Name:        "Carepoint",
		CashBalance: 10.5,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       pharmacyUID[:],
		ProductID: productUID[:],
		Name:      "True Barrier (green) (3 per pack)",
		Price:     10.5,
	}))
	suite.NoError(suite.client.Create(suite.ctx, entity.User{
		UID:         uid[:],
		Name:        "TesterTopTransactionUser",
		CashBalance: 100,
	}))

	suite.NoError(suite.purchaseHistoryClient.Create(suite.ctx, entity.PurchaseHistory{
		UID:               uid[:],
		PharmacyUID:       pharmacyUID[:],
		ProductID:         productUID[:],
		TransactionAmount: 21,
		TransactionDate:   startTime.Add(16 * time.Hour),
	}))
	suite.NoError(suite.purchaseHistoryClient.Create(suite.ctx, entity.PurchaseHistory{
		UID:               uid[:],
		PharmacyUID:       pharmacyUID[:],
		ProductID:         productUID[:],
		TransactionAmount: 42,
		TransactionDate:   startTime.Add(48 * time.Hour),
	}))
	suite.NoError(suite.purchaseHistoryClient.Create(suite.ctx, entity.PurchaseHistory{
		UID:               uid[:],
		PharmacyUID:       pharmacyUID[:],
		ProductID:         productUID[:],
		TransactionAmount: 10.5,
		TransactionDate:   endTime.Add(16 * time.Hour),
	}))
	testCases := []struct {
		Label                         string
		TopNumber, StartTime, EndTime int64
		Want                          want
	}{
		{
			Label:     "ListTopTransactionAmountWithTop1AndStartTime20221001ToEndTime20221010ConditionShouldResponseTransactionAmountIs63",
			TopNumber: 1,
			StartTime: startTime.UnixNano() / int64(time.Millisecond),
			EndTime:   endTime.UnixNano() / int64(time.Millisecond),
			Want:      want{},
		},
	}

	for _, tc := range testCases {
		result, err := suite.client.ListTopTransactionAmount(suite.ctx, tc.TopNumber, tc.StartTime, tc.EndTime)
		suite.NoError(err)
		suite.Equal(uid[:], result.TopTransactionAmountUsers[0].UID)
		suite.Equal("TesterTopTransactionUser", result.TopTransactionAmountUsers[0].Name)
		suite.Equal(float64(63), result.TopTransactionAmountUsers[0].TransactionAmount)
	}
}

func (suite *UserSuite) TestGetTransactionTotalMethod() {
	type want struct {
		Error error
	}

	uid, err := uuid.NewUUID()
	suite.NoError(err)
	pharmacyUID, err := uuid.NewUUID()
	suite.NoError(err)
	productUID, err := uuid.NewUUID()
	suite.NoError(err)

	startTime, err := time.Parse("2006-01-02", "2022-09-01")
	suite.NoError(err)
	endTime, err := time.Parse("2006-01-02", "2022-09-10")
	suite.NoError(err)

	suite.NoError(suite.pharmacyClient.Create(suite.ctx, entity.Pharmacy{
		UID:         pharmacyUID[:],
		Name:        "Carepoint",
		CashBalance: 10.5,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       pharmacyUID[:],
		ProductID: productUID[:],
		Name:      "True Barrier (green) (3 per pack)",
		Price:     10.5,
	}))
	suite.NoError(suite.client.Create(suite.ctx, entity.User{
		UID:         uid[:],
		Name:        "TesterGetTransactionTotal",
		CashBalance: 100,
	}))

	suite.NoError(suite.purchaseHistoryClient.Create(suite.ctx, entity.PurchaseHistory{
		UID:               uid[:],
		PharmacyUID:       pharmacyUID[:],
		ProductID:         productUID[:],
		TransactionAmount: 21,
		TransactionDate:   startTime.Add(16 * time.Hour),
	}))
	suite.NoError(suite.purchaseHistoryClient.Create(suite.ctx, entity.PurchaseHistory{
		UID:               uid[:],
		PharmacyUID:       pharmacyUID[:],
		ProductID:         productUID[:],
		TransactionAmount: 42,
		TransactionDate:   startTime.Add(48 * time.Hour),
	}))
	suite.NoError(suite.purchaseHistoryClient.Create(suite.ctx, entity.PurchaseHistory{
		UID:               uid[:],
		PharmacyUID:       pharmacyUID[:],
		ProductID:         productUID[:],
		TransactionAmount: 10.5,
		TransactionDate:   endTime.Add(16 * time.Hour),
	}))
	testCases := []struct {
		Label              string
		StartTime, EndTime int64
		Want               want
	}{
		{
			Label:     "GetTransactionTotalWithStartTime20220901ToEndTime20220910ConditionShouldResponseTotalIs1TransactionAmountIs63",
			StartTime: startTime.UnixNano() / int64(time.Millisecond),
			EndTime:   endTime.UnixNano() / int64(time.Millisecond),
			Want:      want{},
		},
	}

	for _, tc := range testCases {
		result, err := suite.client.GetTransactionTotal(suite.ctx, tc.StartTime, tc.EndTime)
		suite.NoError(err)
		suite.Equal(int64(1), result.Total)
		suite.Equal(float64(63), result.TransactionAmount)
	}
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
