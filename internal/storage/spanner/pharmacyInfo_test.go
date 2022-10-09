package spanner

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"math"
	"phantom_mask/internal/entity"
	"testing"
	"time"
)

type PharmacyInfoSuite struct {
	suite.Suite
	ctx        context.Context
	logger     *zap.Logger
	client     *PharmacyInfo
	pharmacyID []byte
}

func (suite *PharmacyInfoSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, err := zap.NewDevelopment()
	suite.NoError(err)
	suite.logger = logger
	suite.client = NewPharmacyInfo(suite.logger, session)

	pharmacyClient := NewPharmacy(suite.logger, session)
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	suite.NoError(pharmacyClient.Create(suite.ctx, entity.Pharmacy{
		UID:         uid[:],
		Name:        "TesterPharmacy",
		CashBalance: 10.5,
	}))
	suite.pharmacyID = uid[:]
}

func (suite *PharmacyInfoSuite) TestCreateMethod() {
	type want struct {
		Error error
	}
	var timeNow = time.Now
	defer gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	h, err := time.ParseDuration(fmt.Sprintf("%dh%dm", timeNow().Hour(), timeNow().Minute()))
	suite.NoError(err)

	testCases := []struct {
		Label  string
		Entity entity.PharmacyInfo
		Want   want
	}{
		{
			Label: "CreatePharmacyInfoShouldSuccess",
			Entity: entity.PharmacyInfo{
				UID:       suite.pharmacyID,
				Day:       int64(timeNow().Weekday()),
				OpenHour:  math.Round(h.Hours()*100) / 100,
				CloseHour: math.Round(h.Hours()*100) / 100,
			},
			Want: want{},
		},
		{
			Label: "CreateDuplicatePharmacyShouldFail",
			Entity: entity.PharmacyInfo{
				UID:       suite.pharmacyID,
				Day:       int64(timeNow().Weekday()),
				OpenHour:  math.Round(h.Hours()*100) / 100,
				CloseHour: math.Round(h.Hours()*100) / 100,
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

func TestPharmacyInfoSuite(t *testing.T) {
	suite.Run(t, new(PharmacyInfoSuite))
}
