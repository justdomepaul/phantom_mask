package entity

import (
	"github.com/justdomepaul/toolbox/entity"
	"time"
)

// PRIMARY KEY(UID)
type User struct {
	UID         []byte    `spanner:"UID" json:"uid,omitempty" validate:"required,max=16"`
	Name        string    `spanner:"Name" json:"name,omitempty" validate:"required"`
	CashBalance float64   `spanner:"CashBalance" json:"cash_balance,omitempty" validate:"required"`
	CreatedTime time.Time `spanner:"CreatedTime" json:"created_time,omitempty"`
}

type UserList struct {
	entity.CommonListResponse
	Users []*User `spanner:"Users" json:"users,omitempty"`
}

type TopTransactionAmountUser struct {
	UID               []byte  `spanner:"UID" json:"uid,omitempty"`
	Name              string  `spanner:"Name" json:"name,omitempty"`
	TransactionAmount float64 `spanner:"TransactionAmount" json:"transaction_amount,omitempty"`
}

type TopTransactionAmountList struct {
	TopTransactionAmountUsers []*TopTransactionAmountUser `spanner:"TopTransactionAmountUsers" json:"top_transaction_amount_users,omitempty"`
}

type TopTransactionAmountUserJSON struct {
	*TopTransactionAmountUser
	UID string `json:"uid,omitempty"`
}

type TopTransactionAmountListJSON struct {
	TopTransactionAmountUsers []*TopTransactionAmountUserJSON `json:"top_transaction_amount_users,omitempty"`
}

type TransactionTotal struct {
	Total             int64   `spanner:"Total" json:"total,omitempty"`
	TransactionAmount float64 `spanner:"TransactionAmount" json:"transaction_amount,omitempty"`
}
