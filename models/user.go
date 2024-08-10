package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `json:"-"`
	Address   string         `json:"address"`
}
