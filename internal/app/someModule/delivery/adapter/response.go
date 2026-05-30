package adapter

import (
	"time"

	"github.com/saibaend/template-svc/internal/app/someModule/model"
)

type ItemResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ItemListResponse struct {
	Items  []ItemResponse `json:"items"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ToItemResponse(item model.Item) ItemResponse {
	return ItemResponse{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func ToItemListResponse(items []model.Item, limit, offset int) ItemListResponse {
	resp := ItemListResponse{
		Items:  make([]ItemResponse, 0, len(items)),
		Limit:  limit,
		Offset: offset,
	}

	for _, item := range items {
		resp.Items = append(resp.Items, ToItemResponse(item))
	}

	return resp
}
