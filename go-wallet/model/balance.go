package model

import "time"

type WalletBalance struct {
	AboveThreshold bool       `json:"above_threshold" gorm:"column:above_threshold"`
	WalletId       uint64     `json:"wallet_id"`
	Amount         int64      `json:"amount"`
	LastTopUp      *time.Time `json:"last_topup,omitempty"`
}
