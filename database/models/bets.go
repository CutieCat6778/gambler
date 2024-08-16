package models

import (
	"gambler/backend/database/models/customTypes"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Bet struct {
	gorm.Model
	Name        string                `json:"name" gorm:"unique"`
	Description string                `json:"description"`
	UserBets    []UserBet             `json:"user_bets" gorm:"foreignKey:BetID"`
	BetOptions  pq.StringArray        `json:"betOptions" gorm:"type:text[]"`
	Status      customTypes.BetStatus `json:"status"`
}

type UserBet struct {
	gorm.Model
	UserID    string `json:"user"`
	BetID     uint   `json:"bet_id"` // Foreign key field
	Amount    int    `json:"amount" gorm:"not null"`
	BetOption string `json:"bet_option"`
}
