package storage

import (
	"context"
	"phantom_mask/internal/entity"
)

type ProductEnumType int

const (
	ProductSpecifyPharmacy ProductEnumType = iota
)

type ProductListCondition struct {
	Fields     []ProductEnumType
	PharmacyID []byte
}

func WithProductSpecifyPharmacy(condition ProductListCondition, pharmacyID []byte) ProductListCondition {
	condition.Fields = append(condition.Fields, ProductSpecifyPharmacy)
	condition.PharmacyID = pharmacyID
	return condition
}

type IProduct interface {
	Create(ctx context.Context, input entity.Product) error
	Purchase(ctx context.Context, userID, pharmacyID, productID []byte, quantity int) error
	// List method
	// row required, and min is 1
	// page required, and min is 1
	List(ctx context.Context, row, page uint64, orderEnum OrderListEnum, condition ProductListCondition) (*entity.ProductList, error)
}
