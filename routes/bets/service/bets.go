package service

import (
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/handlers"
	"gambler/backend/handlers/websocket"
	"gambler/backend/tools"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

type (
	CreateBetReq struct {
		Name        string   `json:"name" validate:"required,min=3,max=50,ascii"`
		Description string   `json:"description" validate:"required,min=3,max=50,ascii"`
		BetOptions  []string `json:"betOptions" validate:"required,dive,min=2,max=50,ascii"`
		InputBet    float64  `json:"inputBet" validate:"required,min=1"`
		InputOption string   `json:"inputOption" validate:"required"`
		EndsAt      string   `json:"endsAt" validate:"required"`
	}
	PlaceBetReq struct {
		Amount float64 `json:"amount" validate:"required,min=1"`
		Option string  `json:"option" validate:"required"`
	}
)

func PlaceBet(c *fiber.Ctx) error {

	req := new(PlaceBetReq)

	if err := c.BodyParser(req); err != nil {
		return tools.ReturnData(c, 400, nil, -1)
	}

	if errs := handlers.VHandler.Validate(req); len(errs) > 0 && errs[0].Error {
		return tools.ReturnData(c, 400, errs, -1)
	}

	claims := c.Locals("claims").(jwt.MapClaims)
	userID, jwtErr := claims.GetSubject()
	if jwtErr != nil {
		return tools.ReturnData(c, 400, nil, -1)
	}

	betID := c.Params("id")
	bet, err := handlers.Cache.GetBetById(tools.ParseUInt(betID))
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	user, err := handlers.DB.GetUserByID(tools.ParseUInt(userID))
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	if bet.EndsAt.Before(time.Now()) {
		return tools.ReturnData(c, 400, nil, tools.BET_NOT_ACTIVE)
	}
	if bet.Status != customTypes.Open {
		return tools.ReturnData(c, 400, nil, tools.BET_NOT_ACTIVE)
	}
	if tools.Contains(bet.BetOptions, req.Option) == false {
		log.Info(bet.BetOptions, req.Option)
		return tools.ReturnData(c, 400, nil, tools.BET_OPTION_NOT_FOUND)
	}

	userBet := models.UserBet{
		UserID:    user.ID,
		BetID:     bet.ID,
		Amount:    req.Amount,
		BetOption: req.Option,
	}

	err = handlers.DB.PlaceBet(userBet)
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	err = handlers.Cache.UpdateBet(bet.ID)
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	err = handlers.DB.UpdateUserBalance(-req.Amount, *user, fmt.Sprintf("Placed bet on %s", bet.Name))
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	err = websocket.WebSocket.UpdateBet(bet.ID)
	if err != -1 {
		log.Info("Failed to update bet in websocket")
	}

	return tools.ReturnData(c, 200, true, -1)
}

func GetAllBetsHandler(c *fiber.Ctx) error {
	query := c.QueryInt("type", 0)
	switch query {
	case 0:
		return GetAllActiveBets(c)
	case 1:
		return GetAllPendingBets(c)
	case 2:
		return GetAllClosedBets(c)
	case 3:
		return GetAllCancelledBets(c)
	default:
		return GetAllActiveBets(c)
	}
}

func GetAllActiveBets(c *fiber.Ctx) error {
	bets, err := handlers.Cache.GetAllBet()
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}
	return tools.ReturnData(c, 200, bets, -1)
}

func GetAllPendingBets(c *fiber.Ctx) error {
	res := []models.Bet{}
	bets, err := handlers.Cache.GetAllBet()
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	for _, bet := range *bets {
		if bet.Status == customTypes.Pending {
			res = append(res, bet)
		}
	}

	return tools.ReturnData(c, 200, res, -1)
}

func GetAllClosedBets(c *fiber.Ctx) error {
	res := []models.Bet{}
	bets, err := handlers.Cache.GetAllBet()
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	for _, bet := range *bets {
		if bet.Status == customTypes.Closed {
			res = append(res, bet)
		}
	}

	return tools.ReturnData(c, 200, res, -1)
}

func GetAllCancelledBets(c *fiber.Ctx) error {
	res := []models.Bet{}
	bets, err := handlers.Cache.GetAllBet()
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	for _, bet := range *bets {
		if bet.Status == customTypes.Cancelled {
			res = append(res, bet)
		}
	}

	return tools.ReturnData(c, 200, res, -1)
}

func CreateBet(c *fiber.Ctx) error {

	req := new(CreateBetReq)

	if err := c.BodyParser(req); err != nil {
		return tools.ReturnData(c, 400, nil, -1)
	}

	if errs := handlers.VHandler.Validate(req); len(errs) > 0 && errs[0].Error {
		return tools.ReturnData(c, 400, errs, -1)
	}

	userIDString, jwtErr := c.Locals("claims").(jwt.Claims).GetSubject()
	if jwtErr != nil {
		return tools.ReturnData(c, 401, nil, tools.JWT_INVALID)
	}

	userId := tools.ParseUInt(userIDString)

	rand.NewSource(time.Now().UnixNano())

	bet := models.Bet{
		Name:        req.Name,
		Description: req.Description,
		BetOptions:  pq.StringArray(req.BetOptions),
		Status:      customTypes.Open,
		EndsAt:      tools.ParseTimestamp(req.EndsAt),
		Author:      userId,
	}

	log.Info(bet)

	err := handlers.DB.CreateBet(bet, userId, req.InputOption, req.InputBet)
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	websocket.WebSocket.SendMessageToAll([]byte{tools.BET_UPDATE, tools.WEBSOCKET_VERSION, byte(255)})

	return tools.ReturnData(c, 200, bet, -1)
}

func GetBet(c *fiber.Ctx) error {
	paramsId := c.Params("id")

	id := tools.ParseUInt(paramsId)

	bet, err := handlers.DB.GetBetByID(id)
	if err != -1 {
		log.Info(tools.GetErrorString(err))
		return tools.ReturnData(c, 500, nil, err)
	}

	if bet == nil {
		return tools.ReturnData(c, 404, nil, tools.DB_REC_NOTFOUND)
	}

	return tools.ReturnData(c, 200, bet, -1)
}
