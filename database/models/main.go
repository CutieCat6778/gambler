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
	CreatedAt      sql.NullString   `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt      sql.NullString   `json:"updated_at" gorm:"autoUpdateTime:milli"`
	Games          []Games          `json:"games" gorm:"many2many:user_games"`
	BalanceHistory []BalanceHistory `json:"balance_history" gorm:"foreignKey:UserID"`
}

type Games struct {
	gorm.Model
	Type      customTypes.GameType `json:"type"`
	Users     []User               `json:"users" gorm:"many2many:user_games"`
	CreatedAt sql.NullString       `json:"created_at" gorm:"<-:create;autoCreateTime"`
	ClosedAt  sql.NullString       `json:"closed_at" gorm:"autoUpdateTime:milli"`
}

type Bets struct {
	gorm.Model
	UserID      uint           `json:"user_id"`
	Amount      int            `json:"amount"`
	Description string         `json:"description"`
	Status      int            `json:"status"` // 0 - pending, 1 - won, 2 - lost
	BetID       uint           `json:"bet_id"`
	CreatedAt   sql.NullString `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt   sql.NullString `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type BalanceHistory struct {
	gorm.Model
	UserID    uint           `json:"user_id"`
	Amount    int            `json:"amount"`
	Reason    string         `json:"reason"`
	CreatedAt sql.NullString `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt sql.NullString `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
