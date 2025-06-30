package models

import "time"

type Asset struct {
	DealerID    string    `json:"DEALERID"`
	MSISDN      string    `json:"MSISDN"`
	MPIN        string    `json:"MPIN"`
	Balance     float64   `json:"BALANCE"`
	Status      string    `json:"STATUS"`
	TransAmount float64   `json:"TRANSAMOUNT"`
	TransType   string    `json:"TRANSTYPE"`
	Remarks     string    `json:"REMARKS"`
	CreatedAt   time.Time `json:"CREATEDAT"`
	UpdatedAt   time.Time `json:"UPDATEDAT"`
}

type CreateAssetRequest struct {
	DealerID  string  `json:"DEALERID" binding:"required"`
	MSISDN    string  `json:"MSISDN" binding:"required"`
	MPIN      string  `json:"MPIN" binding:"required"`
	Balance   float64 `json:"BALANCE" binding:"required"`
	Status    string  `json:"STATUS" binding:"required"`
	TransType string  `json:"TRANSTYPE" binding:"required"`
	Remarks   string  `json:"REMARKS"`
}

type UpdateAssetRequest struct {
	Balance     float64 `json:"BALANCE" binding:"required"`
	Status      string  `json:"STATUS" binding:"required"`
	TransAmount float64 `json:"TRANSAMOUNT" binding:"required"`
	TransType   string  `json:"TRANSTYPE" binding:"required"`
	Remarks     string  `json:"REMARKS"`
}
