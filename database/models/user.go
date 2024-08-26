package models

import (
	"gorm.io/gorm"
)

type CustomModel struct {
	gorm.Model
}

type User struct {
	CustomModel
	Name                string           `json:"name"`
	Username            string           `json:"username" gorm:"unique"`
	Password            string           `json:"password"`
	Email               string           `json:"email" gorm:"unique"`
	Balance             float64          `json:"balance"`
	BalanceHistory      []BalanceHistory `json:"balance_history" gorm:"foreignKey:UserID"`
	UserBet             []UserBet        `json:"user_bet" gorm:"foreignKey:UserID"`
	RefreshTokenVersion int              `json:"refresh_token_version"`
}

type BalanceHistory struct {
	CustomModel
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
}
