package model

type JSONB map[string]interface{}
type Pagination struct {
	PageSize     int    `json:"pageSize"`
	CurrentPage  int    `json:"currentPage"`
	TotalPages   int    `json:"totalPages"`
	TotalRecords int    `json:"totalRecords"`
	OrderBy      string `json:"orderBy"`
}
