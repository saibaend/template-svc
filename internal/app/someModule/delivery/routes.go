package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/saibaend/template-svc/internal/app/someModule"
)

type handlers struct {
	Usecase someModule.Usecase
}

func AttachRoutes(r *gin.Engine, uc someModule.Usecase) {

	//h := &handlers{
	//	Usecase: uc,
	//}

}
