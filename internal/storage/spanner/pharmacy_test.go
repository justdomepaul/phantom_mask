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
	"time"
)

type PharmacySuite struct {
	suite.Suite
	ctx                context.Context
	logger             *zap.Logger
	client             *Pharmacy
	pharmacyInfoClient *PharmacyInfo
	productClient      *Product
}

func (suite *PharmacySuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, err := zap.NewDevelopment()
	suite.NoError(err)
	suite.logger = logger
	suite.client = NewPharmacy(suite.logger, session)
	suite.pharmacyInfoClient = NewPharmacyInfo(suite.logger, session)
	suite.productClient = NewProduct(suite.logger, session)
}

func (suite *PharmacySuite) TestCreateMethod() {
	type want struct {
		Error error
	}

	uid, err := uuid.NewUUID()
	suite.NoError(err)

	testCases := []struct {
		Label  string
		Entity entity.Pharmacy
		Want   want
	}{
		{
			Label: "CreatePharmacyShouldSuccess",
			Entity: entity.Pharmacy{
				UID:         uid[:],
				Name:        "Carepoint",
				CashBalance: 10.5,
			},
			Want: want{},
		},
		{
			Label: "CreateDuplicatePharmacyShouldFail",
			Entity: entity.Pharmacy{
				UID:         uid[:],
				Name:        "Carepoint",
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

func (suite *PharmacySuite) TestListSpecifyTimeMethod() {
	type want struct {
		Error error
	}
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	suite.NoError(suite.client.Create(suite.ctx, entity.Pharmacy{
		UID:         uid[:],
		Name:        "TesterListSpecifyTime",
		CashBalance: 100,
	}))
	suite.NoError(suite.pharmacyInfoClient.Create(suite.ctx, entity.PharmacyInfo{
		UID:       uid[:],
		Day:       3,
		OpenHour:  20,
		CloseHour: 26,
	}))

	specifyTimestamp := time.Date(2022, 10, 05, 15, 12, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)

	testCases := []struct {
		Label            string
		SpecifyTimestamp int64
		Want             want
	}{
		{
			Label:            "ListSpecifyTimeShouldResponseMeetWithSpecifyTimestamp",
			SpecifyTimestamp: specifyTimestamp,
			Want:             want{},
		},
	}

	for _, tc := range testCases {
		result, err := suite.client.ListSpecifyTime(suite.ctx, 10, 1, tc.SpecifyTimestamp, storage.PharmacyNameASC)
		suite.NoError(err)
		suite.T().Log(result.Pharmacies)
		suite.Equal(uid[:], result.Pharmacies[0].UID)
		suite.Equal("TesterListSpecifyTime", result.Pharmacies[0].Name)
		suite.Equal(float64(100), result.Pharmacies[0].CashBalance)
		suite.Equal(int64(3), result.Pharmacies[0].Day)
		suite.Equal(float64(20), result.Pharmacies[0].OpenHour)
		suite.Equal(float64(2), result.Pharmacies[0].CloseHour)
	}
}

func (suite *PharmacySuite) TestListByProductPriceRangeMethod() {
	type want struct {
		Error error
	}
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	productID, err := uuid.NewUUID()
	suite.NoError(err)
	productID2, err := uuid.NewUUID()
	suite.NoError(err)
	productID3, err := uuid.NewUUID()
	suite.NoError(err)
	suite.NoError(suite.client.Create(suite.ctx, entity.Pharmacy{
		UID:         uid[:],
		Name:        "TesterListSpecifyTime",
		CashBalance: 100,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       uid[:],
		ProductID: productID[:],
		Name:      "TestByRangeProduct",
		Price:     20,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       uid[:],
		ProductID: productID2[:],
		Name:      "TestByRangeProduct2",
		Price:     30,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       uid[:],
		ProductID: productID3[:],
		Name:      "TestByRangeProduct3",
		Price:     70,
	}))

	testCases := []struct {
		Label string
		Min   int64
		Max   int64
		Want  want
	}{
		{
			Label: "ListByProductPriceRangeWithPriceMin20AndMax50ConditionShouldResponseTwoProducts",
			Min:   20,
			Max:   50,
			Want:  want{},
		},
	}

	for _, tc := range testCases {
		condition := storage.PharmacyListCondition{}
		condition = storage.WithPharmacyProductPriceRange(condition, tc.Min, tc.Max)
		result, err := suite.client.ListByProductPriceRange(suite.ctx, 10, 1, storage.CreatedTimeASC, condition)
		suite.NoError(err)
		suite.Equal(uid[:], result.Pharmacies[0].UID)
		suite.Equal("TesterListSpecifyTime", result.Pharmacies[0].Name)
		suite.Equal(float64(100), result.Pharmacies[0].CashBalance)
	}
}

func (suite *PharmacySuite) TestListPharmacyMixProductMethod() {
	type want struct {
		Error error
	}
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	productID, err := uuid.NewUUID()
	suite.NoError(err)
	productID2, err := uuid.NewUUID()
	suite.NoError(err)
	productID3, err := uuid.NewUUID()
	suite.NoError(err)
	suite.NoError(suite.client.Create(suite.ctx, entity.Pharmacy{
		UID:         uid[:],
		Name:        "TesterListPharmacyMixProduct",
		CashBalance: 100,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       uid[:],
		ProductID: productID[:],
		Name:      "TestByMixProductSalt",
		Price:     20,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       uid[:],
		ProductID: productID2[:],
		Name:      "TestByMixProductSalt2",
		Price:     30,
	}))
	suite.NoError(suite.productClient.Create(suite.ctx, entity.Product{
		UID:       uid[:],
		ProductID: productID3[:],
		Name:      "TestByMix3",
		Price:     70,
	}))

	testCases := []struct {
		Label string
		Name  string
		Want  want
	}{
		{
			Label: "ListPharmacyMixProductWithNameConditionShouldResponse",
			Name:  "Salt",
			Want:  want{},
		},
	}

	for _, tc := range testCases {
		result, err := suite.client.ListPharmacyMixProduct(suite.ctx, 10, 1, tc.Name, storage.PharmacyProduct)
		suite.NoError(err)
		suite.T().Log(result.PharmacyProducts)
	}
}

func TestPharmacySuite(t *testing.T) {
	suite.Run(t, new(PharmacySuite))
}
