package routes

import (
	"gses2.app/api/handlers"
	"net/http"
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

func newRoute(path, method string, handler *handlers.RequestHandler) route {
	return route{
		path:    path,
		method:  method,
		handler: handler,
	}
}
