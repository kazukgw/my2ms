package main

import (
	"time"
)

type User struct {
	UserID    string `gorm:"primary_key"`
	Code      int    `gorm:"index:idx_code_name"`
	Name      string `gorm:"index:idx_code_name"`
	Email     string `gorm:"unique;index:idx_email"`
	CreatedAt time.Time
}

func (User) TableName() string {
	return "users"
}
