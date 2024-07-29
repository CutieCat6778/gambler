package tools

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	GlobalErrorHandlerResp struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Body    interface{} `json:"body,omitempty"`
	}
)

const (
	// DATABASE ERRORS
	DB_UNKOWN_ERR = iota
	DB_REC_NOTFOUND
	DB_DUP_KEY

	// JWT ERRORS
	JWT_FAILED_TO_SIGN
	JWT_FAILED_TO_DECODE
	JWT_INVALID
	JWT_EXPIRED
)

var (
	DATABASE    string
	JWT_SECRET  []byte
	HASH_SECRET string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	DATABASE = os.Getenv("POSTGRES_DB")
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	HASH_SECRET = os.Getenv("HASH_SECRET")
	fmt.Println("[ENV] Loaded Enviroment Variables")
	fmt.Println(DATABASE)
}

func ParseUInt(s string) uint {
	var n uint
	fmt.Sscanf(s, "%d", &n)
	return n
}
