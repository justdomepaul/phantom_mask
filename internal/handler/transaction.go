package handler

import (
	"encoding/json"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/timestamp"
	"github.com/justdomepaul/toolbox/utils"
	"go.uber.org/zap"
	"net/http"
	"phantom_mask/internal/entity"
	spannerDB "phantom_mask/internal/storage/spanner"
	"strconv"
)

var (
	GetNowTimestamp = timestamp.GetNowTimestamp
)

func NewTransaction(
	logger *zap.Logger,
	db spannerDB.Set,
) (*Transaction, error) {
	return &Transaction{
		logger: logger,
		db:     db,
	}, nil
}

type Transaction struct {
	logger *zap.Logger
	db     spannerDB.Set
}

func (h *Transaction) BindRoute(route *gin.Engine) {
	adminGroup := route.Group("/transaction")
	{
		v1Group := adminGroup.Group("/v1")
		v1Group.POST("/purchase", h.Purchase)
		v1Group.GET("/transaction/top", h.ListTransactionTop)
		v1Group.GET("/transaction/product", h.GetTransactionTotal)
	}
}

// Process a user purchases a mask from a pharmacy, and handle all relevant data changes in an atomic transaction.
func (h *Transaction) Purchase(c *gin.Context) {
	req := struct {
		UserID     string `json:"user_id,omitempty" validate:"required"`
		PharmacyID string `json:"pharmacy_id,omitempty" validate:"required"`
		ProductID  string `json:"product_id,omitempty" validate:"required"`
		Quantity   int    `json:"quantity,omitempty" validate:"required,min=1"`
	}{}
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		panic(errorhandler.NewErrJSONUnmarshal(err))
	}
	if err := validator.New().Struct(&req); err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	if err := h.db.Product.Purchase(c, utils.ParseUUID(req.UserID), utils.ParseUUID(req.PharmacyID), utils.ParseUUID(req.ProductID), req.Quantity); err != nil {
		panic(errorhandler.NewErrGRPCExecute(err))
	}
	c.String(http.StatusOK, "ok")
}

// The top x users by total transaction amount of masks within a date range.
func (h *Transaction) ListTransactionTop(c *gin.Context) {
	beforeParseTopNumber := c.DefaultQuery("top_number", "10")
	topNumber, err := strconv.ParseInt(beforeParseTopNumber, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseStartTime := c.DefaultQuery("utc0_millisecond_start_timestamp", "0")
	startTime, err := strconv.ParseInt(beforeParseStartTime, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseEndTime := c.DefaultQuery("utc0_millisecond_end_timestamp", strconv.FormatInt(GetNowTimestamp(), 10))
	endTime, err := strconv.ParseInt(beforeParseEndTime, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	result, err := h.db.User.ListTopTransactionAmount(c, topNumber, startTime, endTime)
	if errors.Is(err, errorhandler.ErrInvalidArguments) {
		panic(errorhandler.NewErrVariable(err))
	}
	if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	resp := &entity.TopTransactionAmountListJSON{}
	var topTransactionAmountUser []*entity.TopTransactionAmountUserJSON
	for _, item := range result.TopTransactionAmountUsers {
		topTransactionAmountUser = append(topTransactionAmountUser, &entity.TopTransactionAmountUserJSON{
			TopTransactionAmountUser: item,
			UID:                      utils.FromUUID(item.UID),
		})
	}
	resp.TopTransactionAmountUsers = topTransactionAmountUser
	c.JSON(http.StatusOK, resp)
}

// The total number of masks and dollar value of transactions within a date range.
func (h *Transaction) GetTransactionTotal(c *gin.Context) {
	beforeParseStartTime := c.DefaultQuery("utc0_millisecond_start_timestamp", "0")
	startTime, err := strconv.ParseInt(beforeParseStartTime, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseEndTime := c.DefaultQuery("utc0_millisecond_end_timestamp", strconv.FormatInt(GetNowTimestamp(), 10))
	endTime, err := strconv.ParseInt(beforeParseEndTime, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	result, err := h.db.User.GetTransactionTotal(c, startTime, endTime)
	c.JSON(http.StatusOK, result)
}
