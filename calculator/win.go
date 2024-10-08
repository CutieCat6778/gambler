package calculator

import (
	"fmt"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"math"

	"github.com/gofiber/fiber/v2/log"
)

type (
	BetLog struct {
		BetAmount float64 `json:"amount"`
		BetOption string  `json:"bet_option"`
	}
)

func CalculateWinningAmount(betID uint, userID uint, inputIndex int, userBetted float64) (float64, int) {
	bet, err := handlers.Cache.GetBetById(betID)
	if err != -1 {
		return 0, err
	}
	if bet.Status != customTypes.Open {
		return 0, tools.BET_NOT_ACTIVE
	}
	if inputIndex >= len(bet.BetOptions) {
		return 0, tools.BET_OPTION_NOT_FOUND
	}
	input := bet.BetOptions[inputIndex]

	amount := userBetted // Total bet amount
	sumBet := userBetted // My total bet in that option
	otherWin := 0.0      // Other bets in that option

	for _, bet := range bet.UserBets {
		amount += bet.Amount
		if bet.UserID != userID && bet.BetOption == input {
			otherWin += bet.Amount
		} else if bet.UserID == userID && bet.BetOption == input {
			sumBet += bet.Amount
		}
	}

	if sumBet == 0.0 {
		return 0, -1
	}

	log.Info(fmt.Sprintf("Amount: %v, SumBet: %v, OtherWin: %v", amount, sumBet, otherWin))

	var winAmount float64 // Total amount will win
	if otherWin == 0.0 {
		winAmount = amount - sumBet // If no one bet on that option
	} else {
		winAmount = amount * (sumBet / otherWin) // Total amount beted / ()
	}

	winPercentage := math.Trunc((winAmount+sumBet)/sumBet*100) / 100

	fmt.Println("Winning percentage: ", winPercentage, winAmount, amount, otherWin, sumBet)

	return winPercentage, -1
}

func CalculateWinForExistedBet(betID uint, userID uint, inputIndex int) (float64, int) {
	bet, err := handlers.Cache.GetBetById(betID)
	if err != -1 {
		return 0, err
	}
	if bet.Status != customTypes.Open {
		return 0, tools.BET_NOT_ACTIVE
	}
	input := bet.BetOptions[inputIndex]

	amount := 0.0   // Total Amount will win
	sumBet := 0.0   // User Bet amount in that option
	otherWin := 0.0 // Other's bet amount in that option

	for _, bet := range bet.UserBets {
		amount += bet.Amount
		if bet.BetOption == input && bet.UserID != userID {
			otherWin += bet.Amount
		} else if bet.BetOption == input && bet.UserID == userID {
			sumBet += bet.Amount
		}
	}

	if sumBet == 0.0 {
		return 0, -1
	}

	var winAmount float64
	if otherWin == 0.0 {
		winAmount = amount
	} else {
		winAmount = amount * (sumBet / otherWin)
	}

	winPercentage := math.Trunc(winAmount/sumBet*100) / 100

	fmt.Println("Winning percentage: ", winPercentage, winAmount, amount, otherWin, sumBet)

	return winPercentage, -1
}
