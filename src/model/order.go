package model

type OrderGetAll struct {
	Page  int
	Count int
}

type OrderInsertRequest struct {
	CancellationReason string `json:"cancellationReason"`
	Note               string `json:"note"`
	OrderDetail        JSONB  `json:"orderDetail"`
	UserID             string `json:"userID"`
	StatusID           string `json:"statusID"`
	ServiceID          string `json:"serviceID"`
	Deadline           string `json:"deadline"`
}
