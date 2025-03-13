package model

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

// type JSONB json.RawMessage

type OrderInsertRequest struct {
	CancellationReason string `json:"cancellationReason"`
	Note               string `json:"note"`
	OrderDetail        JSONB  `json:"orderDetail"`
	UserID             string `json:"userID"`
	StatusID           string `json:"statusID"`
	ServiceID          string `json:"serviceID"`
	Deadline           string `json:"deadline"`
	OrderTypeID        string `json:"orderTypeID"`
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
	OrderID string `json:"orderID"`
	Price   string `json:"price"`
}
