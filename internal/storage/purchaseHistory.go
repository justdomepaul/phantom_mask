package storage

import (
	"context"
	"phantom_mask/internal/entity"
)

type PurchaseHistoryEnumType int

const (
	PurchaseHistoryRefreshTokenNotNull PurchaseHistoryEnumType = iota
)

type PurchaseHistoryListCondition struct {
	Fields []PurchaseHistoryEnumType
}

type IPurchaseHistory interface {
	Create(ctx context.Context, input entity.PurchaseHistory) error
}
