package customTypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// Define the BetStatus type and its constants
type BetStatus string

const (
	Open      BetStatus = "Open"
	Pending   BetStatus = "Pending"
	Closed    BetStatus = "Closed"
	Cancelled BetStatus = "Cancelled"
)

// Implement the sql.Scanner interface for BetStatus
func (bs *BetStatus) Scan(value interface{}) error {
	// Try to assert the value to an int64 (the default type for numeric database columns)
	val, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to scan BetStatus value:", value))
	}

	// Cast the int64 to BetStatus and assign it
	*bs = BetStatus(val)
	return nil
}

// Implement the driver.Valuer interface for BetStatus
func (bs BetStatus) Value() (driver.Value, error) {
	// Return the int representation of the BetStatus
	return string(bs), nil
}
