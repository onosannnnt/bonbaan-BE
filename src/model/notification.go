package model

type NotificationInsertRequest struct {
	Header  string `json:"header"`
	Body    string `json:"body"`
	UserID  string `json:"userID"`
	OrderID string `json:"orderID"`
}
