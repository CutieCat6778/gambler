package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name"`
	Username  string         `json:"username" gorm:"unique"`
	Password  string         `json:"password"`
	Email     string         `json:"email" gorm:"unique"`
	Balance   int            `json:"balance"`
	CreatedAt sql.NullString `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt sql.NullString `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
