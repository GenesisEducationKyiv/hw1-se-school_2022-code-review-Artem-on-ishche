package http

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func getRouterAndRecorder() (*gin.Engine, *httptest.ResponseRecorder) {
	return getRouter(), getRecorder()
}

func getRouter() *gin.Engine {
	router := gin.Default()

	router.Handle(testBtcToUahHandler.GetMethod(), testBtcToUahHandler.GetPath(), testBtcToUahHandler.HandleRequest)
	router.Handle(testSubscribeRequestHandler.GetMethod(), testSubscribeRequestHandler.GetPath(), testSubscribeRequestHandler.HandleRequest)
	router.Handle(testSendEmailsHandler.GetMethod(), testSendEmailsHandler.GetPath(), testSendEmailsHandler.HandleRequest)

	return router
}

func getRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
