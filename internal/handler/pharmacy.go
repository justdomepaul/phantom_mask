package handler

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/utils"
	"go.uber.org/zap"
	"math"
	"net/http"
	"phantom_mask/internal/entity"
	"phantom_mask/internal/storage"
	spannerDB "phantom_mask/internal/storage/spanner"
	"strconv"
)

var MaxInt64Str = strconv.FormatInt(math.MaxInt64, 10)

func NewPharmacy(
	logger *zap.Logger,
	db spannerDB.Set,
) (*Pharmacy, error) {
	return &Pharmacy{
		logger: logger,
		db:     db,
	}, nil
}

type Pharmacy struct {
	logger *zap.Logger
	db     spannerDB.Set
}

func (h *Pharmacy) BindRoute(route *gin.Engine) {
	adminGroup := route.Group("/pharmacy")
	{
		v1Group := adminGroup.Group("/v1")
		v1Group.GET("/", h.ListPharmacy)
		v1Group.GET("/mix", h.ListMix)
		v1Group.GET("/:PharmacyID/product", h.ListProduct)
		v1Group.GET("/product/price", h.ListByProductPriceRange)
	}
}

// ListPharmacy :List all pharmacies open at a specific time and on a day of the week if requested.
func (h *Pharmacy) ListPharmacy(c *gin.Context) {
	beforeParsePage := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(beforeParsePage, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseRow := c.DefaultQuery("row", "10")
	row, err := strconv.ParseUint(beforeParseRow, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	beforeParseSpecifyTimestamp := c.DefaultQuery("specify_utc0_millisecond_timestamp", "0")
	specifyTimestamp, err := strconv.ParseInt(beforeParseSpecifyTimestamp, 10, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	result, err := h.db.Pharmacy.ListSpecifyTime(c, row, page, specifyTimestamp, storage.PharmacyNameASC)
	if errors.Is(err, errorhandler.ErrInvalidArguments) {
		panic(errorhandler.NewErrVariable(err))
	}
	if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	resp := &entity.PharmacySpecifyListJSON{
		CommonListResponse: result.CommonListResponse,
	}
	var pharmacies []*entity.PharmacySpecifyItemJSON
	for _, item := range result.Pharmacies {
		pharmacies = append(pharmacies, &entity.PharmacySpecifyItemJSON{
			PharmacySpecifyTimestamp: item,
			UID:                      utils.FromUUID(item.UID),
		})
	}
	resp.Pharmacies = pharmacies
	c.JSON(http.StatusOK, resp)
}

// Search for pharmacies or masks by name, ranked by relevance to the search term.
func (h *Pharmacy) ListMix(c *gin.Context) {
	beforeParsePage := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(beforeParsePage, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseRow := c.DefaultQuery("row", "10")
	row, err := strconv.ParseUint(beforeParseRow, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	result, err := h.db.Pharmacy.ListPharmacyMixProduct(c, row, page, c.Query("name"), storage.PharmacyProduct)
	if errors.Is(err, errorhandler.ErrInvalidArguments) {
		panic(errorhandler.NewErrVariable(err))
	}
	if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	resp := &entity.PharmacyProductListJSON{
		CommonListResponse: result.CommonListResponse,
	}
	var pharmacyProducts []*entity.PharmacyProductJSON
	for _, item := range result.PharmacyProducts {
		pharmacyProducts = append(pharmacyProducts, &entity.PharmacyProductJSON{
			PharmacyProduct: item,
			UID:             utils.FromUUID(item.UID),
			ProductID:       utils.FromUUID(item.ProductID),
		})
	}
	resp.PharmacyProducts = pharmacyProducts
	c.JSON(http.StatusOK, resp)
}

// List all masks sold by a given pharmacy, sorted by mask name or price.
func (h *Pharmacy) ListProduct(c *gin.Context) {
	valid := struct {
		PharmacyID string `validate:"required"`
	}{}
	beforeParsePage := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(beforeParsePage, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseRow := c.DefaultQuery("row", "10")
	row, err := strconv.ParseUint(beforeParseRow, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	sorted := c.DefaultQuery("sorted", "name")

	pharmacyID := c.Param("PharmacyID")
	valid.PharmacyID = pharmacyID
	if err := validator.New().Struct(&valid); err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	var order storage.OrderListEnum
	switch sorted {
	case "name":
		order = storage.ProductNameASC
	case "price":
		order = storage.ProductPriceASC
	}

	condition := storage.ProductListCondition{}
	condition = storage.WithProductSpecifyPharmacy(condition, utils.ParseUUID(pharmacyID))
	result, err := h.db.Product.List(c, row, page, order, condition)
	if errors.Is(err, errorhandler.ErrInvalidArguments) {
		panic(errorhandler.NewErrVariable(err))
	}
	if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	resp := &entity.ProductListJSON{
		CommonListResponse: result.CommonListResponse,
	}
	var products []*entity.ProductItemJSON
	for _, item := range result.Products {
		products = append(products, &entity.ProductItemJSON{
			Product:   item,
			UID:       utils.FromUUID(item.UID),
			ProductID: utils.FromUUID(item.ProductID),
		})
	}
	resp.Products = products
	c.JSON(http.StatusOK, resp)
}

// List all pharmacies with more or less than x mask products within a price range.
func (h *Pharmacy) ListByProductPriceRange(c *gin.Context) {
	beforeParsePage := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(beforeParsePage, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseRow := c.DefaultQuery("row", "10")
	row, err := strconv.ParseUint(beforeParseRow, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	beforeParseMin := c.DefaultQuery("min", "0")
	min, err := strconv.ParseInt(beforeParseMin, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	beforeParseMax := c.DefaultQuery("max", MaxInt64Str)
	max, err := strconv.ParseInt(beforeParseMax, 0, 64)
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}

	condition := storage.PharmacyListCondition{}
	condition = storage.WithPharmacyProductPriceRange(condition, min, max)
	result, err := h.db.Pharmacy.ListByProductPriceRange(c, row, page, storage.PharmacyNameASC, condition)
	if errors.Is(err, errorhandler.ErrInvalidArguments) {
		panic(errorhandler.NewErrVariable(err))
	}
	if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	resp := &entity.PharmacyListJSON{
		CommonListResponse: result.CommonListResponse,
	}
	var pharmacies []*entity.PharmacyItemJSON
	for _, item := range result.Pharmacies {
		pharmacies = append(pharmacies, &entity.PharmacyItemJSON{
			Pharmacy: item,
			UID:      utils.FromUUID(item.UID),
		})
	}
	resp.Pharmacies = pharmacies
	c.JSON(http.StatusOK, resp)
}
