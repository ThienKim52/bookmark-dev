package dto

// BookmarkCreateInput is the request body for creating a bookmark.
type BookmarkCreateInput struct {
    Description string `json:"description" binding:"required" example:"Google"`
    URL         string `json:"url" binding:"required,url" example:"https://www.google.com"`
}