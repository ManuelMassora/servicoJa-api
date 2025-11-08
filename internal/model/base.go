package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel provides a common ID type (int64) and timestamps for all models.
type BaseModel struct {
    ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
