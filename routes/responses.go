package routes

import (
	"log"
	"net/http"
)

type httpResponseSender interface {
	sendResponse(responseWriter http.ResponseWriter, responseCode int, message string)
}

type httpResponseSenderImpl struct{}

func (sender httpResponseSenderImpl) sendResponse(responseWriter http.ResponseWriter, statusCode int, message string) {
	responseWriter.WriteHeader(statusCode)

	_, err := responseWriter.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error when sending an http response with status code %v and message %v: %v", statusCode, message, err)
	}
}

var responseSender httpResponseSender = httpResponseSenderImpl{}
