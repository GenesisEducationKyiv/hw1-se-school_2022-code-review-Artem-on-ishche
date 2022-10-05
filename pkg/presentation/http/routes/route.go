package routes

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/presentation/http/handlers"
)

type RequestRoute interface {
	GetPath() string
	GetMethod() string
	ProcessRequest(ctx *gin.Context)
}

type handlerFunction func(ctx *gin.Context) *handlers.JSONResponse

func sendJSONResponse(ctx *gin.Context, handle handlerFunction) {
	response := handle(ctx)

	ctx.JSON(response.Code, response.Data)
}
