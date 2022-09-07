package entity

import (
	"github.com/Calmantara/go-common/entity"
	"github.com/Calmantara/go-common/topic"
)

type EmitterPayload struct {
	Id      uint64             `json:"id" gorm:"not null;primaryKey;unique;type:serial;column:id"`
	Issuer  string             `json:"issuer" gorm:"not null"`
	Message string             `json:"message" gorm:"not null"`
	Topic   topic.EmitterTopic `json:"topic" gorm:"not null"`
	Status  bool               `json:"status" gorm:"default:false"`
	Attempt int                `json:"attempt"`
	entity.DefaultColumn
}
