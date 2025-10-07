package user

import "context"

type Repository interface {
	Create(ctx context.Context, name, email, passwordHash string) (User, error)
	ByID(ctx context.Context, id int) (User, error)
	UpdatePartial(ctx context.Context, id int, name, email *string, passwordHash *string) (User, error)
	Delete(ctx context.Context, id int) error
}
