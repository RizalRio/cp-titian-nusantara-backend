package models

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}