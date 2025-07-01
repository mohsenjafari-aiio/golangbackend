package port

import (
	"context"

	"github.com/mohsenjafari-aiio/aiiobackend/internal/order/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, o *domain.Order) error
	GetByID(ctx context.Context, id int64) (*domain.Order, error)
}
