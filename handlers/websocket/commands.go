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
	switch event {
	case tools.BET_ACTION:
		// Handle bet event
		res, err = betEventHandler(data, uuid)
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

func betEventHandler(data []byte, uuid string) ([]byte, int) {
	betType := int(data[0]) // 0: Bet, 1: Cancel
	betID := int(data[1])
	input := int(data[2])
	amountInt := int(data[3])
	amountFrac := int(data[4])
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

	if input >= len(bet.BetOptions) || input < 0 {
		return []byte{}, tools.BET_OPTION_NOT_FOUND
	}

	option := bet.BetOptions[input]

	if betType == 0 {
		err := _handlePlaceBet(bet, option, amount, user)
		if err != -1 {
			return []byte{}, err
		}
	} else if betType == 1 {
		err := _handleCancelBet(bet, option, user)
		if err != -1 {
			return []byte{}, err
		}
	}

	return []byte{tools.BET_ACTION_RES, tools.WEBSOCKET_VERSION, 1}, -1
}

func _handlePlaceBet(bet *models.Bet, input string, amount float64, user *models.User) int {
	if bet.Status != customTypes.Open {
		return tools.BET_NOT_ACTIVE
	}

	userBet := models.UserBet{
		UserID:    user.ID,
		BetID:     bet.ID,
		Amount:    amount,
		BetOption: input,
	}

	err := handlers.DB.PlaceBet(userBet)
	if err != -1 {
		return err
	}

	// Update user balance
	err = handlers.DB.UpdateUserBalance(-amount, *user, fmt.Sprintf("Bet on: %s", bet.Name))
	if err != -1 {
		return err
	}

	err = handlers.Cache.UpdateBet(bet.ID)
	if err != -1 {
		return err
	}

	return -1
}

func _handleCancelBet(bet *models.Bet, input string, user *models.User) int {
	if bet.Status != customTypes.Open {
		return tools.BET_NOT_ACTIVE
	}

	userBet, err := handlers.DB.GetUserBetByBetID(bet.ID, user.Username)
	if err != -1 {
		return err
	}
	err = handlers.DB.CancelBet(*userBet, *user)
	if err != -1 {
		return err
	}

	// Update user balance
	err = handlers.DB.UpdateUserBalance(userBet.Amount, *user, fmt.Sprintf("Cancel bet on: %s", bet.Name))
	if err != -1 {
		return err
	}

	err = handlers.Cache.UpdateBet(bet.ID)
	if err != -1 {
		return err
	}

	return -1
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
