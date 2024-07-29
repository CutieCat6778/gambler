package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name"`
	Username  string         `json:"username" gorm:"unique"`
	Password  []byte         `json:"password"`
	Email     string         `json:"email"`
	Balance   int            `json:"balance"`
	CreatedAt sql.NullString `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt sql.NullString `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
