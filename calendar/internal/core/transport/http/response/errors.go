package response

// ErrorResponse represents the standard JSON body for HTTP error responses.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
