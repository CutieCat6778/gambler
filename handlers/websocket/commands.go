package websocket

import (
	"encoding/json"
	"fmt"
	"gambler/backend/calculator"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"math"

	"github.com/gofiber/fiber/v2/log"
)

func HandleMessageEvent(wsh *WebSocketHandler, uuid string, event int, data []byte) {
	var res []byte
	var err int
	var resp = false
	log.Info("Handling message event:", event, uuid)
	switch event {
	case tools.BET_INFO:
		// Handle bet info event
		res, err = betInfoEventHandler(data, uuid)
		resp = true
	case tools.PING:
		// Handle ping event
		res = []byte{tools.PONG, tools.WEBSOCKET_VERSION}
		err = -1
		resp = true
	default:
		res, err = nil, tools.WS_COMMAND_NOTFOUND
		resp = true
	}

	if err != -1 {
		wsh.SendErrorMessage(uuid, err, tools.GetErrorString(err))
	}

	var wsErr int
	if resp {
		wsErr = wsh.SendMessageToUser(uuid, res)
	}
	log.Info(fmt.Sprintf("%v", res))
	if wsErr != -1 {
		wsh.SendErrorMessage(uuid, wsErr, tools.GetErrorString(wsErr))
	}
}

func betInfoEventHandler(data []byte, uuid string) ([]byte, int) {
	var betLog []calculator.BetLog
	betID := data[0]
	input := int(data[1])
	log.Info(int(betID))
	jsonErr := json.Unmarshal(data[2:], &betLog)
	if jsonErr != nil {
		log.Info(jsonErr)
		return []byte{}, tools.JSON_UNMARSHAL_ERROR
	}

	log.Info(betID, input)
	user, err := handlers.DB.GetUserByID(tools.ParseUInt(uuid))
	if err != -1 {
		return []byte{}, err
	}

	// Calculate winning amount
	winAmount, err := calculator.CalculateWinningAmount(uint(betID), user.ID, input, betLog)
	if err != -1 {
		log.Info(err, tools.GetErrorString(err))
		return []byte{}, err
	}

	intPart, fracPart := math.Modf(winAmount)

	betIDChunks := tools.ChunkBigNumber(int(betID))
	intPartChunks := tools.ChunkBigNumber(int(intPart))
	fracPartChunks := tools.ChunkBigNumber(int(fracPart * 100))

	result := []byte{tools.BET_INFO_RES, tools.WEBSOCKET_VERSION, byte(len(betIDChunks)), byte(len(intPartChunks)), byte(len(fracPartChunks))}
	result = append(result, betIDChunks...)
	result = append(result, intPartChunks...)
	result = append(result, fracPartChunks...)
	return result, -1
}

func combineToFloat64(before, after int) float64 {
	// Convert 'before' directly to float64
	floatBefore := float64(before)

	// Calculate the fractional part by dividing 'after' by the appropriate power of 10
	numDigits := len(fmt.Sprintf("%d", after)) // Count digits in 'after'
	floatAfter := float64(after) / math.Pow(10, float64(numDigits))

	// Combine the two parts
	combined := floatBefore + floatAfter

	return combined
}
