package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Jurusan struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Nama      string    `gorm:"size:100;uniqueIndex;not null" json:"nama"`
	Kode      string    `gorm:"size:20;uniqueIndex;not null" json:"kode"`
	CreatedAt time.Time `json:"created_at"`
}

func (j *Jurusan) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
