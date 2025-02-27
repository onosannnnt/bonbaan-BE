package model

// type JSONB json.RawMessage

type OrderInsertRequest struct {
	CancellationReason string `json:"cancellationReason"`
	Note               string `json:"note"`
	OrderDetail        JSONB  `json:"orderDetail"`
	UserID             string `json:"userID"`
	StatusID           string `json:"statusID"`
	ServiceID          string `json:"serviceID"`
	Deadline           string `json:"deadline"`
}

type ChargeEvent struct {
	Data struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}
