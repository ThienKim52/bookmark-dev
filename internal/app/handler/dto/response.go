package dto

import "github.com/ThienKim52/golang-dev/internal/app/model"

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type BookmarkInput struct {
	Page  int `form:"page" example:"1"`
	Limit int `form:"limit" example:"10"`
}

type BookmarkListResponse struct {
	Data       []*model.Bookmark `json:"data"`
	Pagination *Pagination       `json:"pagination"`
}

type BookmarkCreateResponse struct {
	Data       *model.Bookmark `json:"data"`
	Message		string      `json:"message"`
}

