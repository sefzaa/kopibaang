package dto

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(code int, message string, data interface{}) BaseResponse {
	return BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code int, message string, errors interface{}) BaseResponse {
	return BaseResponse{
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}