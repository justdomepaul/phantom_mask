//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/database/spanner"
	"github.com/justdomepaul/toolbox/jwt"
	"github.com/justdomepaul/toolbox/restful"
	"github.com/justdomepaul/toolbox/stringtool"
	zapLogger "github.com/justdomepaul/toolbox/zap"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"phantom_mask/internal/handler"
	"phantom_mask/internal/storage"
	spannerDB "phantom_mask/internal/storage/spanner"
)

func ctx() context.Context {
	return context.Background()
}

var ctxSet = wire.NewSet(ctx)

var LoggerSet = wire.NewSet(zapLogger.NewLogger)

type Empty struct{}

func RunRestfulServer(logger *zap.Logger, coreOptions config.Set, route *gin.Engine, commonHandler restful.CommonHandler, handlers handler.Set) (Empty, func(), error) {
	handler.AddRoutes(route, commonHandler, handlers)
	pprof.Register(route)

	h2s := &http2.Server{}

	httpServer := &http.Server{
		Addr:    stringtool.StringJoin(":" + coreOptions.Server.Port),
		Handler: h2c.NewHandler(route, h2s),
	}

	go func(s *http.Server) {
		logger.Info("start restful server",
			zap.String("system", coreOptions.Core.SystemName),
			zap.String("port", coreOptions.Server.Port),
		)
		if err := s.ListenAndServe(); err != nil {
			logger.Warn("restful server error or closed",
				zap.String("system", coreOptions.Core.SystemName),
				zap.Error(err),
			)
		}
	}(httpServer)

	return Empty{}, func() {
		ctx, cancel := context.WithTimeout(context.Background(), coreOptions.Server.ServerTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Warn("restful server Failed to Shutdown",
				zap.String("system", coreOptions.Core.SystemName),
				zap.Error(err),
			)
		}
	}, nil
}

func RestfulRunner() (Empty, func(), error) {
	panic(wire.Build(wire.NewSet(
		ctxSet,
		wire.NewSet(
			config.NewSet,
			config.NewCore,
			config.NewServer,
			config.NewSpanner,
			config.NewJWT,
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
		wire.NewSet(jwt.NewEHS384JWTFromOptions, wire.Bind(new(jwt.IJWT), new(*jwt.EHS384JWT))),
		wire.NewSet(restful.NewBasicGuardValidator, wire.Bind(new(restful.GuarderValidator), new(*restful.BasicGuardValidator))),
		wire.NewSet(restful.NewJWTGuarder),
		wire.NewSet(restful.NewGin),
		wire.Value(restful.CommonHandler{
			Error404:   restful.Error404Set,
			QuickReply: restful.QuickReplySet,
			PromHTTP:   restful.NewPromHTTPSet,
		}),
		wire.NewSet(
			handler.NewPharmacy,
			handler.NewTransaction,
			wire.Struct(new(handler.Set), "*")),
		wire.NewSet(restful.NewRender),
		RunRestfulServer,
	)))
}
