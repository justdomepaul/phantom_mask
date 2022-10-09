package entity

import (
	"github.com/justdomepaul/toolbox/entity"
	"time"
)

// Product KEY(UID)
type Product struct {
	UID         []byte    `spanner:"UID" json:"uid,omitempty" validate:"required,max=16"`
	ProductID   []byte    `spanner:"ProductID" json:"product_id,omitempty" validate:"required,max=16"`
	Name        string    `spanner:"Name" json:"name,omitempty" validate:"required"`
	Price       float64   `spanner:"Price" json:"price,omitempty" validate:"required"`
	CreatedTime time.Time `spanner:"CreatedTime" json:"created_time,omitempty"`
}

type ProductList struct {
	entity.CommonListResponse
	Products []*Product `spanner:"Products" json:"products,omitempty"`
}

type ProductItemJSON struct {
	*Product
	UID       string `json:"uid,omitempty"`
	ProductID string `json:"product_id,omitempty"`
}

type ProductListJSON struct {
	entity.CommonListResponse
	Products []*ProductItemJSON `json:"products,omitempty"`
}
