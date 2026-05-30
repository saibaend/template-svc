package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/saibaend/template-svc/internal/app/someModule"
)

func AttachRoutes(r *gin.Engine, uc someModule.Usecase) {
	h := &handlers{usecase: uc}

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		v1.POST("/items", h.CreateItem)
		v1.GET("/items", h.ListItems)
		v1.GET("/items/:id", h.GetItem)
		v1.PUT("/items/:id", h.UpdateItem)
		v1.DELETE("/items/:id", h.DeleteItem)
	}
}
