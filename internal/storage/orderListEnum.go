package storage

type OrderListEnum int

const (
	CreatedTimeASC OrderListEnum = iota
	CreatedTimeDESC
	UpdatedTimeASC
	UpdatedTimeDESC
	PharmacyNameASC
	ProductNameASC
	ProductPriceASC
	PharmacyProduct
)

type OrderPrimaryKeyListEnum bool

const (
	PrimaryKeyASC  OrderPrimaryKeyListEnum = true
	PrimaryKeyDESC                         = false
)
