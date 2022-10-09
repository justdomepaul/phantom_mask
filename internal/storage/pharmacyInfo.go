package storage

import (
	"context"
	"phantom_mask/internal/entity"
)

type IPharmacyInfo interface {
	Create(ctx context.Context, input entity.PharmacyInfo) error
}
