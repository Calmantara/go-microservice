package entity

import "github.com/Calmantara/go-common/entity"

type Wallet struct {
	Remark string `json:"remark" `
	Id     uint64 `json:"id,omitempty" gorm:"not null;primaryKey;unique;type:serial;column:id"`
	UserId uint64 `json:"user_id,omitempty"`
	entity.DefaultColumn
}
