package port

import (
	"context"

	"github.com/mohsenjafari-aiio/aiiobackend/internal/user/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Save(ctx context.Context, u *domain.User) error
}
