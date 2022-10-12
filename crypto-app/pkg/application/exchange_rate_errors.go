package application

import "errors"

var (
	ErrAPIRequestUnsuccessful     = errors.New("API request has been unsuccessful")
	ErrAPIResponseUnmarshallError = errors.New("error when unmarshalling API response")
)
