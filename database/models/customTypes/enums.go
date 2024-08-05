package customTypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// Define the GameType type and its constants
type GameType int

const (
	BlackJack GameType = iota + 1
	Roulette
	Slots
)

// Implement the sql.Scanner interface for GameType
func (gt *GameType) Scan(value interface{}) error {
	// Try to assert the value to an int64 (the default type for numeric database columns)
	val, ok := value.(int64)
	if !ok {
		return errors.New(fmt.Sprint("Failed to scan GameType value:", value))
	}

	// Cast the int64 to GameType and assign it
	*gt = GameType(val)
	return nil
}

// Implement the driver.Valuer interface for GameType
func (gt GameType) Value() (driver.Value, error) {
	// Return the int representation of the GameType
	return int64(gt), nil
}
