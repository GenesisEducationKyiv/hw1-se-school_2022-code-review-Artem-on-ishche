package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/presentation/http/handlers"
)

type RequestRoute interface {
	GetPath() string
	GetMethod() string
	ProcessRequest(ctx *gin.Context)
}

type handlerFunction func(ctx *gin.Context) *handlers.JSONResponse

func sendJSONResponse(logger services.Logger, ctx *gin.Context, handle handlerFunction) {
	response := handle(ctx)

	ctx.JSON(response.Code, response.Data)
	logger.Info("My response code: " + strconv.Itoa(response.Code))
}
