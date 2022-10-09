package storage

import (
	"context"
	"phantom_mask/internal/entity"
)

type PharmacyEnumType int

const (
	PharmacyProductPriceRange PharmacyEnumType = iota
)

type PharmacyListCondition struct {
	Fields []PharmacyEnumType
	Min    int64
	Max    int64
}

func WithPharmacyProductPriceRange(condition PharmacyListCondition, min, max int64) PharmacyListCondition {
	condition.Fields = append(condition.Fields, PharmacyProductPriceRange)
	condition.Min = min
	condition.Max = max
	return condition
}

type IPharmacy interface {
	Create(ctx context.Context, input entity.Pharmacy) error
	// ListPharmacyMixProduct method
	// row required, and min is 1
	// page required, and min is 1
	ListPharmacyMixProduct(ctx context.Context, row, page uint64, name string, orderEnum OrderListEnum) (*entity.PharmacyProductList, error)
	// ListSpecifyTime method
	// row required, and min is 1
	// page required, and min is 1
	ListSpecifyTime(ctx context.Context, row, page uint64, specifyTimestamp int64, orderEnum OrderListEnum) (*entity.PharmacySpecifyTimestampList, error)
	// ListByProductPriceRange method
	// row required, and min is 1
	// page required, and min is 1
	ListByProductPriceRange(ctx context.Context, row, page uint64, orderEnum OrderListEnum, condition PharmacyListCondition) (*entity.PharmacyList, error)
}
