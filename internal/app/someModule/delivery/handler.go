package delivery

import (
	"errors"
	"github.com/saibaend/template-svc/internal/app/someModule/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saibaend/template-svc/internal/app/someModule"
	"github.com/saibaend/template-svc/internal/app/someModule/delivery/adapter"
)

type handlers struct {
	usecase someModule.Usecase
}

// CreateItem godoc
//
//	@Summary		Create item
//	@Description	Create a new item
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			body	body		adapter.CreateItemRequest	true	"Item payload"
//	@Success		201		{object}	adapter.ItemResponse
//	@Failure		400		{object}	adapter.ErrorResponse
//	@Failure		500		{object}	adapter.ErrorResponse
//	@Router			/v1/items [post]
func (h *handlers) CreateItem(c *gin.Context) {
	var req adapter.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	item, err := h.usecase.Create(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, adapter.ToItemResponse(*item))
}

// GetItem godoc
//
//	@Summary		Get item by ID
//	@Description	Returns a single item
//	@Tags			items
//	@Produce		json
//	@Param			id	path		int	true	"Item ID"
//	@Success		200	{object}	adapter.ItemResponse
//	@Failure		400	{object}	adapter.ErrorResponse
//	@Failure		404	{object}	adapter.ErrorResponse
//	@Failure		500	{object}	adapter.ErrorResponse
//	@Router			/v1/items/{id} [get]
func (h *handlers) GetItem(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, adapter.ToItemResponse(*item))
}

// ListItems godoc
//
//	@Summary		List items
//	@Description	Returns paginated items
//	@Tags			items
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"	default(20)
//	@Param			offset	query		int	false	"Offset"	default(0)
//	@Success		200		{object}	adapter.ItemListResponse
//	@Failure		400		{object}	adapter.ErrorResponse
//	@Failure		500		{object}	adapter.ErrorResponse
//	@Router			/v1/items [get]
func (h *handlers) ListItems(c *gin.Context) {
	limit, err := parseQueryInt(c, "limit", 20)
	if err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	offset, err := parseQueryInt(c, "offset", 0)
	if err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	items, err := h.usecase.List(c.Request.Context(), limit, offset)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, adapter.ToItemListResponse(items, limit, offset))
}

// UpdateItem godoc
//
//	@Summary		Update item
//	@Description	Updates item by ID
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Item ID"
//	@Param			body	body		adapter.UpdateItemRequest	true	"Item payload"
//	@Success		200		{object}	adapter.ItemResponse
//	@Failure		400		{object}	adapter.ErrorResponse
//	@Failure		404		{object}	adapter.ErrorResponse
//	@Failure		500		{object}	adapter.ErrorResponse
//	@Router			/v1/items/{id} [put]
func (h *handlers) UpdateItem(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	var req adapter.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	item, err := h.usecase.Update(c.Request.Context(), id, req.Title, req.Description)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, adapter.ToItemResponse(*item))
}

// DeleteItem godoc
//
//	@Summary		Delete item
//	@Description	Deletes item by ID
//	@Tags			items
//	@Param			id	path	int	true	"Item ID"
//	@Success		204
//	@Failure		400	{object}	adapter.ErrorResponse
//	@Failure		404	{object}	adapter.ErrorResponse
//	@Failure		500	{object}	adapter.ErrorResponse
//	@Router			/v1/items/{id} [delete]
func (h *handlers) DeleteItem(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid item id")
	}

	return id, nil
}

func parseQueryInt(c *gin.Context, key string, defaultValue int) (int, error) {
	raw := c.Query(key)
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.New("invalid " + key)
	}

	return value, nil
}

func writeUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		writeError(c, http.StatusNotFound, "not_found", err.Error())
	case err.Error() == "title is required", err.Error() == "invalid item id":
		writeError(c, http.StatusBadRequest, "bad_request", err.Error())
	default:
		writeError(c, http.StatusInternalServerError, "internal_error", err.Error())
	}
}

func writeError(c *gin.Context, status int, code, message string) {
	c.JSON(status, adapter.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
