package adapter

type CreateItemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateItemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
