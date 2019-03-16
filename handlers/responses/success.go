package responses

// SuccessResponse is a type representing a successful request's json response
type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

// NewSuccessResponse creates a new SuccessResponse
func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{data}
}
