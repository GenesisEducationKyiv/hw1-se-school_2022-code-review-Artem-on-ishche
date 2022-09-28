package http

import (
	"github.com/gin-gonic/gin"
)

type RequestHandler interface {
	GetPath() string
	GetMethod() string
	HandleRequest(*gin.Context)
}
