//go:build wireinject

package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/database/spanner"
	zapTool "github.com/justdomepaul/toolbox/zap"
	"go.uber.org/zap"
	"os"
	"phantom_mask/internal/entity"
	"phantom_mask/internal/storage"
	spannerDB "phantom_mask/internal/storage/spanner"
	"phantom_mask/internal/utils"
	"time"
)

func ctx() context.Context {
	return context.Background()
}

var ctxSet = wire.NewSet(ctx)

var LoggerSet = wire.NewSet(zapTool.NewLogger)

type Empty struct{}

func NewImportData(c context.Context, db spannerDB.Set) *ImportData {
	return &ImportData{
		ctx: c,
		db:  db,
	}
}

type ImportData struct {
	ctx context.Context
	db  spannerDB.Set
}

func (i ImportData) Init() error {
	var users []entity.UserJSON
	var pharmacies []entity.PharmacyJSON

	userFile, err := os.Open("./data/users.json")
	if err != nil {
		return err
	}
	defer userFile.Close()
	if err := json.NewDecoder(userFile).Decode(&users); err != nil {
		return err
	}
	pharmacyFile, err := os.Open("./data/pharmacies.json")
	if err != nil {
		panic(err)
	}
	defer pharmacyFile.Close()
	if err := json.NewDecoder(pharmacyFile).Decode(&pharmacies); err != nil {
		return err
	}

	pharmacyMap := map[string][]byte{}
	productMap := map[string][]byte{}
	// pharmacy data
	for _, phy := range pharmacies {
		phyUID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		if err := i.db.Pharmacy.Create(i.ctx, entity.Pharmacy{
			UID:         phyUID[:],
			Name:        phy.Name,
			CashBalance: phy.CashBalance,
		}); err != nil {
			return err
		}
		pharmacyMap[phy.Name] = phyUID[:]

		openCloseHour, err := utils.ParseTimeFormat(phy.OpeningHours)
		if err != nil {
			return err
		}
		for _, info := range openCloseHour {
			if err := i.db.PharmacyInfo.Create(i.ctx, entity.PharmacyInfo{
				UID:       phyUID[:],
				Day:       info.Day,
				OpenHour:  info.OpenHour,
				CloseHour: info.CloseHour,
			}); err != nil {
				return err
			}
		}

		for _, mask := range phy.Masks {
			productUID, err := uuid.NewUUID()
			if err != nil {
				return err
			}
			if err := i.db.Product.Create(i.ctx, entity.Product{
				UID:       phyUID[:],
				ProductID: productUID[:],
				Name:      mask.Name,
				Price:     mask.Price,
			}); err != nil {
				return err
			}
			productMap[mask.Name] = productUID[:]
		}
	}

	// userData
	for _, us := range users {
		userID, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		if err := i.db.User.Create(i.ctx, entity.User{
			UID:         userID[:],
			Name:        us.Name,
			CashBalance: us.CashBalance,
		}); err != nil {
			return err
		}
		for _, usHis := range us.PurchaseHistories {
			specifyTime, err := time.Parse("2006-01-02 15:04:05", usHis.TransactionDate)
			if err != nil {
				return err
			}
			if err := i.db.PurchaseHistory.Create(i.ctx, entity.PurchaseHistory{
				UID:               userID[:],
				PharmacyUID:       pharmacyMap[usHis.PharmacyName],
				ProductID:         productMap[usHis.MaskName],
				TransactionAmount: usHis.TransactionAmount,
				TransactionDate:   specifyTime,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func Run(logger *zap.Logger, coreOptions config.Set, initData *ImportData) (Empty, func(), error) {
	if err := initData.Init(); err != nil {
		return Empty{}, nil, err
	}
	return Empty{}, func() {}, nil
}

func Runner() (Empty, func(), error) {
	panic(wire.Build(wire.NewSet(
		ctxSet,
		wire.NewSet(
			config.NewSet,
			config.NewCore,
			config.NewGRPC,
			config.NewJWT,
			config.NewServer,
			config.NewSpanner,
		),
		LoggerSet,
		spanner.NewExtendSpannerDatabase,
		wire.NewSet(
			wire.NewSet(spannerDB.NewPharmacy, wire.Bind(new(storage.IPharmacy), new(*spannerDB.Pharmacy))),
			wire.NewSet(spannerDB.NewPharmacyInfo, wire.Bind(new(storage.IPharmacyInfo), new(*spannerDB.PharmacyInfo))),
			wire.NewSet(spannerDB.NewProduct, wire.Bind(new(storage.IProduct), new(*spannerDB.Product))),
			wire.NewSet(spannerDB.NewUser, wire.Bind(new(storage.IUser), new(*spannerDB.User))),
			wire.NewSet(spannerDB.NewPurchaseHistory, wire.Bind(new(storage.IPurchaseHistory), new(*spannerDB.PurchaseHistory))),
			wire.Struct(new(spannerDB.Set), "*")),
		wire.NewSet(NewImportData),
		Run,
	)))
}
