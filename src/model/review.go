package model

type ReviewInsertRequest struct {
	Rating    int    `json:"rating"`
	Detail    string `json:"detail"`
	ServiceID string `json:"serviceID"`
	OrderID   string `json:"orderID"`
}
