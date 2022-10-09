package storage

import (
	"context"
	"phantom_mask/internal/entity"
)

type UserEnumType int

const (
	UserRefreshTokenNotNull UserEnumType = iota
)

type UserListCondition struct {
	Fields []UserEnumType
}

type IUser interface {
	Create(ctx context.Context, input entity.User) error
	ListTopTransactionAmount(ctx context.Context, topNumber, startTime, endTime int64) (*entity.TopTransactionAmountList, error)
	GetTransactionTotal(ctx context.Context, startTime, endTime int64) (*entity.TransactionTotal, error)
}
