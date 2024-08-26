package websocket

import (
	"encoding/json"
	"fmt"
	"gambler/backend/calculator"
	"gambler/backend/tools"
	"math"

	"github.com/gofiber/fiber/v2/log"
)

func HandleMessageEvent(wsh *WebSocketHandler, uuid string, event int, data []byte) {
	var res []byte
	var err int
	switch event {
	case tools.BET:
		// Handle bet event
		res, err = betEventHandler(event, data)
	case tools.BET_INFO:
		// Handle bet info event
		res, err = betInfoEventHandler(event, data)
	default:
		res, err = nil, tools.WS_COMMAND_NOTFOUND
	}

	if err != -1 {
		wsh.SendErrorMessage(uuid, err, tools.GetErrorString(err))
	}

	wsErr := wsh.SendMessageToUser(uuid, res)
	log.Info(fmt.Sprintf("%v", res))
	if wsErr != nil {
		wsh.SendErrorMessage(uuid, tools.WS_UNKNOWN_ERR, tools.GetErrorString(tools.WS_UNKNOWN_ERR))
	}
}

func betEventHandler(event int, data []byte) ([]byte, int) {
	return []byte{}, -1
}

func betInfoEventHandler(event int, data []byte) ([]byte, int) {
	var betLog []calculator.BetLog
	betID := data[0]
	input := int(data[1])
	log.Info(int(betID))
	jsonErr := json.Unmarshal(data[2:], &betLog)
	if jsonErr != nil {
		log.Info(jsonErr)
		return []byte{}, tools.JSON_UNMARSHAL_ERROR
	}

	// Calculate winning amount
	winAmount, err := calculator.CalculateWinningAmount(fmt.Sprintf("b-%d", int(betID)), input, betLog)
	if err != -1 {
		log.Info(err, tools.GetErrorString(err))
		return []byte{}, err
	}

	intPart, fracPart := math.Modf(winAmount)

	return []byte{tools.BET_INFO_RES, tools.WEBSOCKET_VERSION, byte(intPart), byte(int(fracPart * 100))}, -1
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
