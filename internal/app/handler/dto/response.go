package dto

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type BookmarkPaginationInput struct {
	Page  int `form:"page" example:"1"`
	Limit int `form:"limit" example:"10"`
}

type SuccessResponse[Data any] struct {
	Message    string             `json:"message,omitempty"`
	Data       Data 			  `json:"data,omitempty"`
	Pagination *Pagination        `json:"pagination,omitempty"`
}
