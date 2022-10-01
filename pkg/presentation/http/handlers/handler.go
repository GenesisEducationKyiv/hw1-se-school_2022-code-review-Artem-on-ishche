package handlers

import "github.com/gin-gonic/gin"

type RequestHandler interface {
	HandleRequest(ctx *gin.Context) *JSONResponse
}
