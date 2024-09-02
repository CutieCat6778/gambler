package calculator

import (
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"math"

	"github.com/gofiber/fiber/v2/log"
	logger "github.com/gofiber/fiber/v2/log"
)

type (
	BetLog struct {
		BetAmount float64 `json:"amount"`
		BetOption string  `json:"bet_option"`
	}
)

func generateBetLog(bets []models.UserBet) []BetLog {
	var betLogs []BetLog
	for _, bet := range bets {
		betLogs = append(betLogs, BetLog{
			BetAmount: bet.Amount,
			BetOption: bet.BetOption,
		})
	}
	return betLogs
}

func generateTotalBet(userBets []models.UserBet) map[string]float64 {
	total := make(map[string]float64)
	for _, bet := range userBets {
		total[bet.BetOption] += bet.Amount
	}
	return total
}

func CalculateWinningAmount(betID string, inputIndex int, betLog []BetLog) (float64, int) {
	bet, err := handlers.Cache.GetBetById(betID)
	if err != -1 {
		return 0, err
	}
	if bet.Status != customTypes.Open {
		return 0, tools.BET_NOT_ACTIVE
	}
	input := bet.BetOptions[inputIndex]
	totalBet := generateTotalBet(bet.UserBets)

	amount := 0.0
	sumBet := 0.0
	otherWin := 0.0
	for _, log := range betLog {
		logger.Info(fmt.Sprintf("%v", log))
		if log.BetOption != input {
			sumBet += log.BetAmount
		} else {
			amount += log.BetAmount
		}
	}

	for i, bet := range totalBet {
		log.Info(bet)
		if i != input {
			sumBet += bet
		} else {
			otherWin += bet
		}
	}

	if sumBet == 0.0 {
		return 0, -1
	}

	var winAmount float64
	if otherWin == 0.0 {
		winAmount = sumBet
	} else {
		winAmount = sumBet * (amount / otherWin)
	}

	winPercentage := math.Trunc(winAmount/amount*100) / 100

	fmt.Println("Winning percentage: ", winPercentage, winAmount, amount, otherWin, sumBet)

	return winPercentage, -1
}
