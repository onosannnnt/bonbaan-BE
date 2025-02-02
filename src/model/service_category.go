package model

type AddServiceToCategoryRequest struct {
	CategoryID string `json:"category_id"`
	ServiceID  string `json:"service_id"`
}

type RemoveServiceFromCategoryRequest struct {
	CategoryID string `json:"category_id"`
	ServiceID  string `json:"service_id"`
}