package entity

import "github.com/Calmantara/go-common/entity"

type Balance struct {
	Remark   string `json:"remark" `
	Id       uint64 `json:"id,omitempty" gorm:"not null;primaryKey;unique;type:serial;column:id"`
	WalletId uint64 `json:"wallet_id" gorm:"not null;index"`
	Amount   int64  `json:"amount" gorm:"not null"`
	entity.DefaultColumn
	WalletDetail *Wallet `json:"wallet_detail,omitempty" gorm:"references:WalletId;foreignKey:Id"`
}

func (Balance) TableName() string {
	return "h_balances"
}
