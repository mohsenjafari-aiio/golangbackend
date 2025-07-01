package adapter

import (
	"context"

	"github.com/mohsenjafari-aiio/aiiobackend/internal/product/domain"
	"github.com/mohsenjafari-aiio/aiiobackend/internal/product/port"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) port.ProductRepository {
	return &GormProductRepository{db: db}
}

func (r *GormProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *GormProductRepository) Save(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *GormProductRepository) UpdateStock(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Model(p).Update("stock", p.Stock).Error
}
