package entity

import (
	"time"

	"gorm.io/gorm"
)

type RecordFlag string

const (
	RECORD_ACTIVE  RecordFlag = "ACTIVE"
	RECORD_DELETED RecordFlag = "DELETED"
)

type DefaultColumn struct {
	CreatedBy  string         `json:"created_by,omitempty"`
	UpdatedBy  string         `json:"updated_by,omitempty"`
	RecordFlag RecordFlag     `json:"record_flag,omitempty" gorm:"index;default:'ACTIVE'"`
	CreatedAt  *time.Time     `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
