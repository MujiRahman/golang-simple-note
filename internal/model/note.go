package model

import "time"

type Note struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Title     string    `gorm:"size:255;not null"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
