package model

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

// type JSONB json.RawMessage

type OrderInputRequest struct {
	Price              float64  `json:"price"`
	CancellationReason string   `json:"cancellation_reason"`
	Items              []string `json:"items"`
	PackageID          string   `json:"packageID"`
	ServiceID          string   `json:"serviceID"`
	UserID             string   `json:"userID"`
	Deadline           string   `json:"deadline"`
	Note               string   `json:"note"`
	Vow                string   `json:"vow"`
	VowRecordID        string   `json:"vow_record_id"`
	OrderTypeID        string   `json:"order_type_id"`
}

type ChargeEvent struct {
	Data struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}

type SubmitOrderRequest struct {
	OrderID     string                `json:"orderID"`
	Attachments []Entities.Attachment `json:"attachments,omitempty"`
}

type ConfirmOrderRequest struct {
	OrderID string  `json:"orderID"`
	Price   float64 `json:"price"`
}
