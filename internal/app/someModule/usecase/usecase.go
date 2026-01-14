package usecase

import "github.com/saibaend/template-svc/internal/app/someModule"

type usecase struct {
	service someModule.Service
}

func New(service someModule.Service) usecase {
	return usecase{
		service: service,
	}
}
