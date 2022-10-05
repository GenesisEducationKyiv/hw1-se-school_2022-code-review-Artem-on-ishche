package handlers

type JSONResponse struct {
	Code int
	Data any
}

func NewJSONResponse(code int, data any) *JSONResponse {
	return &JSONResponse{Code: code, Data: data}
}
