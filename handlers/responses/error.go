package responses

// ErrInternalError is a generic error used for all internal errors
var ErrInternalError = Error{
	Title:  "Internal Error",
	Detail: "Encountered internal error, please try again in a few minutes.",
}

// ErrInvalidJSONBody is a generic Error returned when request body could not be parsed as JSON
var ErrInvalidJSONBody = Error{
	Title: "Invalid Body", Detail: "Failed to parse request body as json",
}

// ErrorResponse is a type representing the error response json
// returned to the client
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error is a type used
type Error struct {
	// A short, human-readable summary of the problem that
	// SHOULD NOT change from occurrence to occurrence of the problem
	Title string `json:"title"`

	// A human-readable explanation specific to this occurrence of the problem
	Detail string `json:"detail"`
}

// NewErrorResponse creates a new ErrorResponse from errors
func NewErrorResponse(errors ...Error) ErrorResponse {
	return ErrorResponse{errors}
}
