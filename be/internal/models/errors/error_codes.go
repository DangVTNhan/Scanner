package errors

// Error codes for the application
const (
	// General error codes (1000-1999)
	ErrCodeUnknown           = "ERR1000" // Unknown error
	ErrCodeInvalidRequest    = "ERR1001" // Invalid request
	ErrCodeInvalidParameters = "ERR1002" // Invalid parameters
	ErrCodeNotFound          = "ERR1003" // Resource not found
	ErrCodeUnauthorized      = "ERR1004" // Unauthorized access
	ErrCodeForbidden         = "ERR1005" // Forbidden access
	ErrCodeTimeout           = "ERR1006" // Request timeout
	ErrCodeTooManyRequests   = "ERR1007" // Too many requests
	ErrCodeServerError       = "ERR1008" // Internal server error

	// Database error codes (2000-2999)
	ErrCodeDatabaseConnection = "ERR2000" // Database connection error
	ErrCodeDatabaseQuery      = "ERR2001" // Database query error
	ErrCodeDatabaseInsert     = "ERR2002" // Database insert error
	ErrCodeDatabaseUpdate     = "ERR2003" // Database update error
	ErrCodeDatabaseDelete     = "ERR2004" // Database delete error
	ErrCodeDatabaseNotFound   = "ERR2005" // Database record not found

	// Weather service error codes (3000-3999)
	ErrCodeWeatherServiceConnection = "ERR3000" // Weather service connection error
	ErrCodeWeatherServiceResponse   = "ERR3001" // Weather service response error
	ErrCodeWeatherServiceTimeout    = "ERR3002" // Weather service timeout
	ErrCodeWeatherDataNotAvailable  = "ERR3003" // Weather data not available

	// Report error codes (4000-4999)
	ErrCodeReportGeneration = "ERR4000" // Report generation error
	ErrCodeReportNotFound   = "ERR4001" // Report not found
	ErrCodeReportComparison = "ERR4002" // Report comparison error
	ErrCodeReportInvalid    = "ERR4003" // Invalid report data
)

// ErrorCodeToHTTPStatus maps error codes to HTTP status codes
var ErrorCodeToHTTPStatus = map[string]int{
	// General error codes
	ErrCodeUnknown:           500,
	ErrCodeInvalidRequest:    400,
	ErrCodeInvalidParameters: 400,
	ErrCodeNotFound:          404,
	ErrCodeUnauthorized:      401,
	ErrCodeForbidden:         403,
	ErrCodeTimeout:           408,
	ErrCodeTooManyRequests:   429,
	ErrCodeServerError:       500,

	// Database error codes
	ErrCodeDatabaseConnection: 500,
	ErrCodeDatabaseQuery:      500,
	ErrCodeDatabaseInsert:     500,
	ErrCodeDatabaseUpdate:     500,
	ErrCodeDatabaseDelete:     500,
	ErrCodeDatabaseNotFound:   404,

	// Weather service error codes
	ErrCodeWeatherServiceConnection: 500,
	ErrCodeWeatherServiceResponse:   500,
	ErrCodeWeatherServiceTimeout:    504,
	ErrCodeWeatherDataNotAvailable:  404,

	// Report error codes
	ErrCodeReportGeneration: 500,
	ErrCodeReportNotFound:   404,
	ErrCodeReportComparison: 500,
	ErrCodeReportInvalid:    400,
}
