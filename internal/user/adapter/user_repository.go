package adapter

import (
	"context"

	"github.com/mohsenjafari-aiio/aiiobackend/internal/user/domain"
	"github.com/mohsenjafari-aiio/aiiobackend/internal/user/port"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) port.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Save(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}
