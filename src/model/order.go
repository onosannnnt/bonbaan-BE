package model

type OrderGetAll struct {
	Page  int
	Count int
}

type JSONB map[string]interface{}

type OrderInsertRequest struct {
	CancellationReason string `json:"cancellationReason"`
	Note               string `json:"note"`
	OrderDetail        JSONB  `json:"orderDetail"`
	UserID             string `json:"userID"`
	StatusID           string `json:"statusID"`
	Dateline           string `json:"dateline"`
}
