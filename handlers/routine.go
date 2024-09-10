package handlers

import (
	"gambler/backend/database/models"

	"github.com/gofiber/fiber/v2/log"
)

var (
	Bets  []models.Bet = []models.Bet{}
	Ready int          = 0
)

func CallRoutine(db DBHandler) {
	log.Info("[ROUTINE] Started")

}
