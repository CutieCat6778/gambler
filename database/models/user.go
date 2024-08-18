package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                string           `json:"name"`
	Username            string           `json:"username" gorm:"unique"`
	Password            string           `json:"password"`
	Email               string           `json:"email" gorm:"unique"`
	Balance             int              `json:"balance"`
	BalanceHistory      []BalanceHistory `json:"balance_history" gorm:"foreignKey:UserID"`
	UserBet             []UserBet        `json:"user_bet" gorm:"foreignKey:UserID"`
	RefreshTokenVersion int              `json:"refresh_token_version"`
}

type BalanceHistory struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}
