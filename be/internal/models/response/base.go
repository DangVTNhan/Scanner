package response

// BaseResponse is the standardized response structure for all API responses
type BaseResponse struct {
	Message   string      `json:"message"`   // Human-readable message about the response
	Status    string      `json:"status"`    // Status of the response (success, error)
	ErrorCode string      `json:"errorCode"` // Error code (empty if status is success)
	Data      interface{} `json:"data"`      // Response data (can be any type)
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *BaseResponse {
	return &BaseResponse{
		Message:   message,
		Status:    "success",
		ErrorCode: "",
		Data:      data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, errorCode string, data interface{}) *BaseResponse {
	return &BaseResponse{
		Message:   message,
		Status:    "error",
		ErrorCode: errorCode,
		Data:      data,
	}
}
