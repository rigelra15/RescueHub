package structs

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}