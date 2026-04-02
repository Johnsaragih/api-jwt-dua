package utils

type BaseResponse struct {
	Status  int         `json:"status"`
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	User    interface{} `json:"user,omitempty"`
	Token   string      `json:"access_token,omitempty"`
}
