package adapter

import (
	"context"

	"github.com/mohsenjafari-aiio/aiiobackend/internal/order/domain"
	"github.com/mohsenjafari-aiio/aiiobackend/internal/order/port"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) port.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) Save(ctx context.Context, o *domain.Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *GormOrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).Preload("User").Preload("Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
