package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/toolbox/restful"
)

type Set struct {
	Pharmacy    *Pharmacy
	Transaction *Transaction
}

func AddRoutes(route *gin.Engine, commonHandler restful.CommonHandler, handlers Set) {
	route.GET("/ping", commonHandler.QuickReply)
	route.GET("/metrics", commonHandler.PromHTTP)

	handlers.Pharmacy.BindRoute(route)
	handlers.Transaction.BindRoute(route)

	route.NoRoute(commonHandler.Error404)
}
