package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/saibaend/template-svc/internal/app/someModule"
	"github.com/saibaend/template-svc/internal/app/someModule/model"
)

const (
	defaultListLimit  = 20
	maxListLimit      = 100
	defaultListOffset = 0
)

type usecase struct {
	service someModule.Service
}

func New(service someModule.Service) someModule.Usecase {
	return &usecase{service: service}
}

func (u *usecase) Create(ctx context.Context, title, description string) (*model.Item, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	return u.service.Create(ctx, strings.TrimSpace(title), strings.TrimSpace(description))
}

func (u *usecase) GetByID(ctx context.Context, id int64) (*model.Item, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid item id")
	}

	return u.service.GetByID(ctx, id)
}

func (u *usecase) List(ctx context.Context, limit, offset int) ([]model.Item, error) {
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if offset < 0 {
		offset = defaultListOffset
	}

	return u.service.List(ctx, limit, offset)
}

func (u *usecase) Update(ctx context.Context, id int64, title, description string) (*model.Item, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid item id")
	}
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	return u.service.Update(ctx, id, strings.TrimSpace(title), strings.TrimSpace(description))
}

func (u *usecase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid item id")
	}

	return u.service.Delete(ctx, id)
}

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title is required")
	}

	return nil
}
