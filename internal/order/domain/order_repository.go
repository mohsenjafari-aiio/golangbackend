package domain

import (
	"context"
)

type OrderRepository interface {
	Save(ctx context.Context, o *Order) error
	GetByID(ctx context.Context, id int64) (*Order, error)
}
