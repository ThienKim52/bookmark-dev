package dto


// BookmarkUpdateInput is the request body for updating a bookmark.
type BookmarkUpdateInput struct {
    Description string `json:"description" binding:"required" example:"Google"`
    URL         string `json:"url" binding:"required,url" example:"https://www.google.com"`
}