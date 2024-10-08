package models

import (
	"encoding/json"
	"gambler/backend/database/models/customTypes"
	"time"

	"github.com/lib/pq"
)

type Bet struct {
	CustomModel
	Name        string                `json:"name" gorm:"unique"`
	Description string                `json:"description"`
	UserBets    []UserBet             `json:"user_bets" gorm:"foreignKey:BetID"`
	BetOptions  pq.StringArray        `json:"bet_options" gorm:"type:text[]"`
	Status      customTypes.BetStatus `json:"status"`
	EndsAt      time.Time             `json:"ends_at"`
	Author      uint                  `json:"author,string"`
}

func (b Bet) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Bet) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &b)
}

type UserBet struct {
	CustomModel
	UserID    uint    `json:"user"`
	BetID     uint    `json:"bet_id"` // Foreign key field
	Amount    float64 `json:"amount" gorm:"not null"`
	BetOption string  `json:"bet_option"`
}
