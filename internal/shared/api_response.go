package shared

type SuccessResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ErrorResponse struct {
	Code  int `json:"-"`
	Error any `json:"error"`
}

func NewErrorResponse(code int, err any) *ErrorResponse {
	return &ErrorResponse{
		Code:  code,
		Error: err,
	}
}
