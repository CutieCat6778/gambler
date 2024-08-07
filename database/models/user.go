package models

import (
	"database/sql"
	"gambler/backend/database/models/customTypes"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string           `json:"name"`
	Username       string           `json:"username" gorm:"unique"`
	Password       string           `json:"password"`
	Email          string           `json:"email" gorm:"unique"`
	Balance        int              `json:"balance"`
	Games          []Games          `json:"games" gorm:"many2many:user_games"`
	BalanceHistory []BalanceHistory `json:"balance_history" gorm:"foreignKey:UserID"`
	UserBet        []UserBet        `json:"user_bet" gorm:"foreignKey:UserID"`
}

type Games struct {
	gorm.Model
	Type     customTypes.GameType `json:"type"`
	Users    []User               `json:"users" gorm:"many2many:user_games"`
	ClosedAt sql.NullString       `json:"closed_at" gorm:"autoUpdateTime:milli"`
}

type BalanceHistory struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}
