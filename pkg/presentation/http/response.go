package http

import (
	"log"
	"net/http"
)

type Response struct {
	StatusCode int
	Message    string
}

func newResponse(statusCode int, responseText string) Response {
	return Response{statusCode, responseText}
}

func sendResponse(responseWriter http.ResponseWriter, response Response) {
	responseWriter.WriteHeader(response.StatusCode)

	_, err := responseWriter.Write([]byte(response.Message))
	if err != nil {
		log.Fatalf(
			"Error when sending an http response with status code %v and message %v: %v",
			response.StatusCode, response.Message, err,
		)
	}
}
