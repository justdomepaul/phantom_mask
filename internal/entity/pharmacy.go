package entity

import (
	"github.com/justdomepaul/toolbox/entity"
	"time"
)

// PRIMARY KEY(UID)
type Pharmacy struct {
	UID         []byte    `spanner:"UID" json:"uid,omitempty" validate:"required,max=16"`
	Name        string    `spanner:"Name" json:"name,omitempty" validate:"required"`
	CashBalance float64   `spanner:"CashBalance" json:"cash_balance,omitempty" validate:"required"`
	CreatedTime time.Time `spanner:"CreatedTime" json:"created_time,omitempty"`
}

type PharmacyList struct {
	entity.CommonListResponse
	Pharmacies []*Pharmacy `spanner:"Pharmacies" json:"pharmacies,omitempty"`
}

type PharmacyItemJSON struct {
	*Pharmacy
	UID string `json:"uid,omitempty"`
}

type PharmacyListJSON struct {
	entity.CommonListResponse
	Pharmacies []*PharmacyItemJSON `json:"pharmacies,omitempty"`
}

type PharmacySpecifyTimestamp struct {
	UID         []byte    `spanner:"UID" json:"uid,omitempty"`
	Name        string    `spanner:"Name" json:"name,omitempty"`
	CashBalance float64   `spanner:"CashBalance" json:"cash_balance,omitempty"`
	CreatedTime time.Time `spanner:"CreatedTime" json:"created_time,omitempty"`
	Day         int64     `spanner:"Day" json:"day,omitempty"`
	OpenHour    float64   `spanner:"OpenHour" json:"open_hour,omitempty"`
	CloseHour   float64   `spanner:"CloseHour" json:"close_hour,omitempty"`
}

type PharmacySpecifyTimestampList struct {
	entity.CommonListResponse
	Pharmacies []*PharmacySpecifyTimestamp `spanner:"Pharmacies" json:"pharmacy_specify_timestamp_items,omitempty"`
}

type PharmacySpecifyItemJSON struct {
	*PharmacySpecifyTimestamp
	UID string `spanner:"UID" json:"uid,omitempty"`
}

type PharmacySpecifyListJSON struct {
	entity.CommonListResponse
	Pharmacies []*PharmacySpecifyItemJSON `json:"pharmacies,omitempty"`
}

type PharmacyProduct struct {
	UID          []byte  `spanner:"UID" json:"uid,omitempty"`
	ProductID    []byte  `spanner:"ProductID" json:"product_id,omitempty"`
	PharmacyName string  `spanner:"PharmacyName" json:"pharmacy_name,omitempty"`
	CashBalance  float64 `spanner:"CashBalance" json:"cash_balance,omitempty"`
	ProductName  string  `spanner:"ProductName" json:"product_name,omitempty"`
	Price        float64 `spanner:"Price" json:"price,omitempty"`
}

type PharmacyProductList struct {
	entity.CommonListResponse
	PharmacyProducts []*PharmacyProduct `spanner:"PharmacyProducts" json:"pharmacy_products,omitempty"`
}

type PharmacyProductJSON struct {
	*PharmacyProduct
	UID       string `json:"uid,omitempty"`
	ProductID string `json:"product_id,omitempty"`
}

type PharmacyProductListJSON struct {
	entity.CommonListResponse
	PharmacyProducts []*PharmacyProductJSON `json:"pharmacy_products,omitempty"`
}
