package routes

import (
	"net/http"

	"gses2.app/api/handlers"
)

type route struct {
	path    string
	method  string
	handler *handlers.RequestHandler
}

func (r route) processRequest(responseWriter http.ResponseWriter, request *http.Request) {
	response := (*r.handler).HandleRequest(request)
	responseSender.sendResponse(responseWriter, response.StatusCode, response.Message)
}
