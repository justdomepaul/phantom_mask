package spanner

import (
	"context"
	"github.com/google/uuid"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"phantom_mask/internal/entity"
	"phantom_mask/internal/storage"
	"testing"
)

type ProductSuite struct {
	suite.Suite
	ctx            context.Context
	logger         *zap.Logger
	userClient     *User
	pharmacyClient *Pharmacy
	client         *Product
	userID         []byte
	pharmacyID     []byte
}

func (suite *ProductSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, err := zap.NewDevelopment()
	suite.NoError(err)
	suite.logger = logger
	suite.client = NewProduct(suite.logger, session)
	suite.userClient = NewUser(suite.logger, session)
	suite.pharmacyClient = NewPharmacy(suite.logger, session)
	userID, err := uuid.NewUUID()
	suite.NoError(err)
	pharmacyID, err := uuid.NewUUID()
	suite.NoError(err)
	suite.NoError(suite.userClient.Create(suite.ctx, entity.User{
		UID:         userID[:],
		Name:        "TesterPharmacyUser",
		CashBalance: 3,
	}))
	suite.userID = userID[:]
	suite.NoError(suite.pharmacyClient.Create(suite.ctx, entity.Pharmacy{
		UID:         pharmacyID[:],
		Name:        "TesterPharmacy",
		CashBalance: 3,
	}))
	suite.pharmacyID = pharmacyID[:]
}

func (suite *ProductSuite) TestCreateMethod() {
	type want struct {
		Error error
	}

	uid, err := uuid.NewUUID()
	suite.NoError(err)

	testCases := []struct {
		Label  string
		Entity entity.Product
		Want   want
	}{
		{
			Label: "CreatePharmacyShouldSuccess",
			Entity: entity.Product{
				UID:       suite.pharmacyID,
				ProductID: uid[:],
				Name:      "True Barrier (green) (3 per pack)",
				Price:     10.5,
			},
			Want: want{},
		},
		{
			Label: "CreateDuplicatePharmacyShouldFail",
			Entity: entity.Product{
				UID:       suite.pharmacyID,
				ProductID: uid[:],
				Name:      "True Barrier (green) (3 per pack)",
				Price:     10.5,
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

func (suite *ProductSuite) TestPurchaseMethod() {
	type want struct {
		Error bool
	}
	productID, err := uuid.NewUUID()
	suite.NoError(err)
	suite.NoError(suite.client.Create(suite.ctx, entity.Product{
		UID:       suite.pharmacyID,
		ProductID: productID[:],
		Name:      "TesterPurchaseProductName",
		Price:     0.1,
	}))

	testCases := []struct {
		Label     string
		ProductID []byte
		Quantity  int
		Want      want
	}{
		{
			Label:     "PurchaseProduct5ItemsShouldResponseSuccess",
			ProductID: productID[:],
			Quantity:  5,
			Want:      want{},
		},
		{
			Label:     "PurchaseProduct100ItemsOverUserCashBalanceShouldResponseFail",
			ProductID: productID[:],
			Quantity:  100,
			Want: want{
				Error: true,
			},
		},
	}

	for _, tc := range testCases {
		if tc.Want.Error {
			suite.Error(suite.client.Purchase(suite.ctx, suite.userID, suite.pharmacyID, tc.ProductID, tc.Quantity))
		} else {
			suite.NoError(suite.client.Purchase(suite.ctx, suite.userID, suite.pharmacyID, tc.ProductID, tc.Quantity))
		}
	}
}

func (suite *ProductSuite) TestListMethod() {
	type want struct {
		Error error
	}
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	uid2, err := uuid.NewUUID()
	suite.NoError(err)
	productID, err := uuid.NewUUID()
	suite.NoError(err)
	productID2, err := uuid.NewUUID()
	suite.NoError(err)
	suite.NoError(suite.pharmacyClient.Create(suite.ctx, entity.Pharmacy{
		UID:         uid[:],
		Name:        "TesterListProduct",
		CashBalance: 100,
	}))
	suite.NoError(suite.pharmacyClient.Create(suite.ctx, entity.Pharmacy{
		UID:         uid2[:],
		Name:        "TesterListProduct2",
		CashBalance: 100,
	}))
	suite.NoError(suite.client.Create(suite.ctx, entity.Product{
		UID:       uid[:],
		ProductID: productID[:],
		Name:      "TesterProductName",
		Price:     100,
	}))
	suite.NoError(suite.client.Create(suite.ctx, entity.Product{
		UID:       uid2[:],
		ProductID: productID2[:],
		Name:      "TesterProductName2",
		Price:     100,
	}))

	testCases := []struct {
		Label             string
		SpecifyPharmacyID []byte
		Sorted            storage.OrderListEnum
		Want              want
	}{
		{
			Label:             "ListProductByPharmacyIDAndSortedNameShouldResponseProduct",
			SpecifyPharmacyID: uid[:],
			Sorted:            storage.ProductNameASC,
			Want:              want{},
		},
		{
			Label:             "ListProductByPharmacyIDAndSortedPriceShouldResponseProduct",
			SpecifyPharmacyID: uid[:],
			Sorted:            storage.ProductPriceASC,
			Want:              want{},
		},
	}

	for _, tc := range testCases {
		condition := storage.ProductListCondition{}
		condition = storage.WithProductSpecifyPharmacy(condition, tc.SpecifyPharmacyID)
		result, err := suite.client.List(suite.ctx, 10, 1, tc.Sorted, condition)
		suite.NoError(err)
		suite.Equal(uid[:], result.Products[0].UID)
		suite.Equal(productID[:], result.Products[0].ProductID)
		suite.Equal("TesterProductName", result.Products[0].Name)
		suite.Equal(float64(100), result.Products[0].Price)
	}
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductSuite))
}
