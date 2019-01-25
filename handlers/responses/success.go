package responses

// SuccessResponse is a type representing a successful request's json response
type SuccessResponse struct {
	Success bool `json:"success"`
}

// NewSuccessResponse creates a new SuccessResponse
func NewSuccessResponse() SuccessResponse {
	return SuccessResponse{true}
}
