package websocket

import (
	"encoding/json"
	"fmt"
	"gambler/backend/calculator"
	"gambler/backend/database/models"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"math"

	"github.com/gofiber/fiber/v2/log"
)

func HandleMessageEvent(wsh *WebSocketHandler, uuid string, event int, data []byte) {
	var res []byte
	var err int
	log.Info("Handling message event:", event, uuid)
	switch event {
	case tools.BET_ACTION_BET:
		// Handle bet event
		res, err = betActionBetEventHandler(wsh, data, uuid)
	case tools.BET_ACTION_CANCEL:
		// Handle cancel bet event
		res, err = betActionCancelEventHandler(wsh, data, uuid)
	case tools.BET_INFO:
		// Handle bet info event
		res, err = betInfoEventHandler(data)
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

func betActionBetEventHandler(wsh *WebSocketHandler, data []byte, uuid string) ([]byte, int) {
	betID := int(data[0])
	input := int(data[1])
	amountInt := int(data[2])
	amountFrac := int(data[3])
	amount := combineToFloat64(amountInt, amountFrac)

	user, err := handlers.DB.GetUserByUsername(uuid)
	if err != -1 {
		return []byte{}, err
	}

	if user.Balance < amount {
		return []byte{}, tools.BET_INSUFFICIENT_BALANCE
	}

	bet, err := handlers.Cache.GetBetById(fmt.Sprintf("b-%d", betID))
	if err != -1 {
		return []byte{}, err
	}

	log.Info(fmt.Sprintf("Bet: %d", len(bet.UserBets)))

	if input >= len(bet.BetOptions) || input < 0 {
		return []byte{}, tools.BET_OPTION_NOT_FOUND
	}

	option := bet.BetOptions[input]

	if bet.Status != customTypes.Open {
		return []byte{}, tools.BET_NOT_ACTIVE
	}

	userBet := models.UserBet{
		UserID:    user.ID,
		BetID:     bet.ID,
		Amount:    amount,
		BetOption: option,
	}

	err = handlers.DB.PlaceBet(userBet)
	if err != -1 {
		return []byte{}, err
	}

	// Update user balance
	err = handlers.DB.UpdateUserBalance(-amount, *user, fmt.Sprintf("Bet on: %s", bet.Name))
	if err != -1 {
		return []byte{}, err
	}

	err = handlers.Cache.UpdateBet(bet.ID)
	if err != -1 {
		return []byte{}, err
	}

	betUpdateEventHandler(wsh, bet.ID)

	return []byte{tools.BET_ACTION_RES, tools.WEBSOCKET_VERSION, 1}, -1
}

func betActionCancelEventHandler(wsh *WebSocketHandler, data []byte, uuid string) ([]byte, int) {
	betID := int(data[0])
	userBetID := uint(data[1])

	user, err := handlers.DB.GetUserByUsername(uuid)
	if err != -1 {
		return []byte{}, err
	}

	bet, err := handlers.Cache.GetBetById(fmt.Sprintf("b-%d", betID))
	if err != -1 {
		return []byte{}, err
	}

	log.Info(fmt.Sprintf("Bet: %d", len(bet.UserBets)))

	if bet.Status != customTypes.Open {
		return []byte{}, tools.BET_NOT_ACTIVE
	}

	userBet, err := handlers.DB.GetUserBetByID(userBetID)
	if err != -1 {
		return []byte{}, err
	}

	err = handlers.DB.CancelBet(*userBet, *user)
	if err != -1 {
		return []byte{}, err
	}

	// Update user balance
	err = handlers.DB.UpdateUserBalance(userBet.Amount, *user, fmt.Sprintf("Cancel bet on: %s", bet.Name))
	if err != -1 {
		return []byte{}, err
	}

	err = handlers.Cache.UpdateBet(bet.ID)
	if err != -1 {
		return []byte{}, err
	}

	betUpdateEventHandler(wsh, bet.ID)

	return []byte{}, -1
}

func betInfoEventHandler(data []byte) ([]byte, int) {
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

	// Calculate winning amount
	winAmount, err := calculator.CalculateWinningAmount(fmt.Sprintf("b-%d", int(betID)), input, betLog)
	if err != -1 {
		log.Info(err, tools.GetErrorString(err))
		return []byte{}, err
	}

	intPart, fracPart := math.Modf(winAmount)

	return []byte{tools.BET_INFO_RES, tools.WEBSOCKET_VERSION, byte(intPart), byte(int(fracPart * 100))}, -1
}

func betUpdateEventHandler(wsh *WebSocketHandler, betID uint) int {
	bet, err := handlers.DB.GetBetByID(betID)
	if err != -1 {
		return err
	}
	log.Info(fmt.Sprintf("Bet: %d", len(bet.UserBets)))
	marshal, jsonErr := json.Marshal(bet)
	if jsonErr != nil {
		return tools.JSON_MARSHAL_ERROR
	}
	// Send bet update to all users
	result := []byte{tools.BET_UPDATE, tools.WEBSOCKET_VERSION, byte(bet.ID)}
	result = append(result, marshal...)
	return wsh.SendMessageToAll(result)
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
