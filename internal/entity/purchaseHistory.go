package entity

import (
	"github.com/justdomepaul/toolbox/entity"
	"time"
)

// PRIMARY KEY(UID, TransactionDate)
type PurchaseHistory struct {
	UID               []byte    `spanner:"UID" json:"uid,omitempty" validate:"required,max=16"`
	PharmacyUID       []byte    `spanner:"PharmacyUID" json:"pharmacy_uid,omitempty" validate:"required"`
	ProductID         []byte    `spanner:"ProductID" json:"product_id,omitempty" validate:"required"`
	TransactionAmount float64   `spanner:"TransactionAmount" json:"transaction_amount,omitempty" validate:"required"`
	TransactionDate   time.Time `spanner:"TransactionDate" json:"transaction_date,omitempty" validate:"required"`
}

type PurchaseHistoryList struct {
	entity.CommonListResponse
	PurchaseHistories []*PurchaseHistory `spanner:"PurchaseHistories" json:"purchase_histories,omitempty"`
}
