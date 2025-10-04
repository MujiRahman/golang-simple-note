package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Username  string    `gorm:"uniqueIndex;size:100" json:"username"`
	Password  string    `json:"-"` // hashed
	Email     string    `gorm:"uniqueIndex;size:100;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
