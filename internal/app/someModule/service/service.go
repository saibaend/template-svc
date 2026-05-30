package service

import (
	"context"

	"github.com/saibaend/template-svc/internal/app/someModule"
	"github.com/saibaend/template-svc/internal/app/someModule/model"
)

type service struct {
	repo someModule.Repository
}

func New(repo someModule.Repository) someModule.Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, title, description string) (*model.Item, error) {
	item := &model.Item{
		Title:       title,
		Description: description,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*model.Item, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context, limit, offset int) ([]model.Item, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id int64, title, description string) (*model.Item, error) {
	item := &model.Item{
		ID:          id,
		Title:       title,
		Description: description,
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	deleted, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if !deleted {
		return model.ErrNotFound
	}

	return nil
}
