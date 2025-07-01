package domain

import "context"

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Save(ctx context.Context, u *User) error
}
