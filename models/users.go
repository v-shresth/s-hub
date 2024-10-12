package models

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	gorm.Model
	FirstName       string          `gorm:"column:first_name" json:"firstName"`
	LastName        string          `gorm:"column:last_name" json:"lastName"`
	Email           string          `gorm:"column:email" json:"email"`
	Password        string          `gorm:"column:password" json:"password"`
	UserConnections UserConnections `gorm:"foreignkey:user_id;references:id" json:"userConnections"`
	Sessions        []Session       `gorm:"foreignkey:user_id;references:id" json:"sessions"`
}

type UserConnections struct {
	gorm.Model
	ConnectionString string `gorm:"column:connection_string" json:"connectionString"`
	UserID           uint   `gorm:"column:user_id" json:"userId"`
}

type Session struct {
	gorm.Model
	UserId    uint      `json:"user_id" gorm:"column:user_id"`
	StartedAt time.Time `json:"started_at" gorm:"column:started_at;default:current_timestamp"`
	EndedAt   time.Time `json:"ended_at" gorm:"column:ended_at;default:null"`
}
