package domain

import "context"

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*Product, error)
	Save(ctx context.Context, p *Product) error
	UpdateStock(ctx context.Context, p *Product) error
}
