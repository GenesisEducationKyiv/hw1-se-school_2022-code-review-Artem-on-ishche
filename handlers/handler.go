package handlers

import "net/http"

type httpResponse struct {
	StatusCode int
	Message    string
}

func newHttpResponse(statusCode int, responseText string) httpResponse {
	return httpResponse{statusCode, responseText}
}

type RequestHandler interface {
	HandleRequest(*http.Request) httpResponse
}
