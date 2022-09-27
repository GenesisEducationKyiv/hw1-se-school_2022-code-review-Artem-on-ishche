package http

import "net/http"

type RequestHandler interface {
	GetPath() string
	GetMethod() string
	GetResponse(*http.Request) Response
}

func GetHandlerFunction(handler RequestHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		sendResponse(writer, handler.GetResponse(request))
	}
}
