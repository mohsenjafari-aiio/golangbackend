package port

import (
	"context"

	"github.com/mohsenjafari-aiio/aiiobackend/internal/product/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Save(ctx context.Context, p *domain.Product) error
	UpdateStock(ctx context.Context, p *domain.Product) error
}
