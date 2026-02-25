package models

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}