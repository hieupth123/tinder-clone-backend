package respond

type Respond struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationResponse struct {
	Limit int         `json:"limit"`
	Page  int         `json:"page"`
	Pages int         `json:"pages"`
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}
