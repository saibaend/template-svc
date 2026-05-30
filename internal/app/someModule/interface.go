package someModule

import (
	"context"

	"github.com/saibaend/template-svc/internal/app/someModule/model"
)

type Repository interface {
	Create(ctx context.Context, item *model.Item) error
	GetByID(ctx context.Context, id int64) (*model.Item, error)
	List(ctx context.Context, limit, offset int) ([]model.Item, error)
	Update(ctx context.Context, item *model.Item) error
	Delete(ctx context.Context, id int64) (bool, error)
}

type Service interface {
	Create(ctx context.Context, title, description string) (*model.Item, error)
	GetByID(ctx context.Context, id int64) (*model.Item, error)
	List(ctx context.Context, limit, offset int) ([]model.Item, error)
	Update(ctx context.Context, id int64, title, description string) (*model.Item, error)
	Delete(ctx context.Context, id int64) error
}

type Usecase interface {
	Create(ctx context.Context, title, description string) (*model.Item, error)
	GetByID(ctx context.Context, id int64) (*model.Item, error)
	List(ctx context.Context, limit, offset int) ([]model.Item, error)
	Update(ctx context.Context, id int64, title, description string) (*model.Item, error)
	Delete(ctx context.Context, id int64) error
}
