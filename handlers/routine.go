package handlers

import (
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/tools"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var (
	Bets  []models.Bet = []models.Bet{}
	Ready int          = 0
)

func CallRoutine(db DBHandler) {
	log.Info("[ROUTINE] Started")
	bets, err := db.GetAllBets()
	if err != -1 {
		log.Error(err)
		tools.SendWebHook(tools.GetErrorString(err))
	}
	for _, bet := range *bets {
		res := checkBetExp(bet, db)
		if res == -1 {
			Bets = append(Bets, bet)
		} else if res != -2 {
			tools.SendWebHook(tools.GetErrorString(res))
			continue
		}
	}
	StartRoutine()
	Cache.LoadDatabaseBets()
}

func checkBetExp(bet models.Bet, db DBHandler) int {
	cachedBet, err := Cache.GetBetById(fmt.Sprintf("b-%d", bet.ID))
	if err != -1 {
		log.Error(err)
		tools.SendWebHook(tools.GetErrorString(err))
		return err
	}
	log.Info(cachedBet)
	if bet.Status != customTypes.Open {
		if cachedBet != nil {
			err := Cache.RemoveBet(bet.ID)
			if err != -1 {
				log.Error(err)
				tools.SendWebHook(tools.GetErrorString(err))
				return err
			}
		}
		return -2
	}
	if bet.EndsAt.Before(time.Now()) {
		newBet := bet
		newBet.Status = customTypes.Pending
		err := db.UpdateBet(newBet)
		if err != -1 {
			log.Error(err)
			tools.SendWebHook(tools.GetErrorString(err))
			return err
		}
		err = Cache.UpdateBet(newBet.ID)
		if err != -1 {
			log.Error(err)
			tools.SendWebHook(tools.GetErrorString(err))
			return err
		}
		return -3
	}
	return -1
}

func removeAt(slice []models.Bet, index int) []models.Bet {
	// Check if index is within bounds
	if index < 0 || index >= len(slice) {
		return slice // Return the original slice if index is out of bounds
	}

	// Remove the element by slicing
	return append(slice[:index], slice[index+1:]...)
}

func StartRoutine() {
	ticker := time.NewTicker(time.Minute * 15)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if Ready != 1 {
					return
				}
				for i, bet := range Bets {
					res := checkBetExp(bet, DB)
					if res == -1 {
						continue
					} else if res != -2 {
						tools.SendWebHook(tools.GetErrorString(res))
						continue
					}
					Bets = removeAt(Bets, i)
				}
				return
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
