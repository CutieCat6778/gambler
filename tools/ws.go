package tools

const (
	WS_CONNECTION = iota
	WS_CLOSE
	WS_ERR
	MESSAGE
	BET_ACTION_BET
	BET_ACTION_CANCEL
	BET_INFO
	BET_INFO_RES
	BET_UPDATE
	USER_UPDATE
	PING
	PONG
)

func ChunkBigNumber(n int) []byte {
	largeNumberBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		largeNumberBytes[7-i] = byte(n >> (i * 8))
	}
	return largeNumberBytes
}
