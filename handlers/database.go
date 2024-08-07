package handlers

import (
	"errors"
	"gambler/backend/database"
	"gambler/backend/database/models"
	"gambler/backend/tools"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type DBHandler struct {
	DB *gorm.DB
}

var (
	DB DBHandler
)

func NewDB() DBHandler {
	db := database.InitDatabase()
	DB = DBHandler{db}
	return DB
}

// User methods

func (h DBHandler) CreateUser(user models.User) int {
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the user
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Errorf("Error creating user: %v", err)
		return dbHandleError(err)
	}

	// Ensure the user ID is populated
	if user.ID == 0 {
		tx.Rollback()
		log.Error("User ID not populated after creation")
		return dbHandleError(errors.New("User ID not populated after creation"))
	}

	// Create the initial balance history entry
	initialBalanceHistory := models.BalanceHistory{
		UserID: user.ID,
		Amount: 0,
		Reason: "Initial balance",
	}

	if err := tx.Create(&initialBalanceHistory).Error; err != nil {
		tx.Rollback()
		log.Errorf("Error creating balance history: %v", err)
		return dbHandleError(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return dbHandleError(err)
	}

	return -1
}

func (h DBHandler) UpdateUser(user models.User) (*models.User, int) {
	res := h.DB.Save(&user)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &user, -1
}

func (h DBHandler) GetUserByID(id uint) (*models.User, int) {
	var user models.User
	res := h.DB.First(&user, id)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &user, -1
}

func (h DBHandler) GetUserByUsername(username string) (*models.User, int) {
	var user models.User
	res := h.DB.Preload("BalanceHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Limit(1)
	}).Preload("Games", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Limit(1)
	}).Where("username = ?", username).First(&user)

	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &user, -1
}

func (h DBHandler) DeleteUserByID(id uint) int {
	res := h.DB.Delete(&models.User{}, id)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

// Game methods

func (h DBHandler) CreateGame(game models.Games) int {
	res := h.DB.Model(&game)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

func (h DBHandler) FindGameByUser(user models.User, gameId string) (*models.Games, int) {
	var game models.Games
	res := h.DB.Model(&user).Where("id in ?", gameId).Association("Games").Find(&game)
	if res != nil {
		return nil, dbHandleError(res)
	}
	return &game, -1
}

func (h DBHandler) FindGameByID(gameId string) (*models.Games, int) {
	var game models.Games
	res := h.DB.First(&game, gameId)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &game, -1
}

func (h DBHandler) CloseGame(gameId string, balances []models.BalanceHistory) int {
	game, err := h.FindGameByID(gameId)
	if err != -1 {
		return err
	}
	res := h.DB.Save(&game)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

// BalanceHistory methods

func (h DBHandler) CreateBalanceHistory(balance models.BalanceHistory) int {
	res := h.DB.Create(&balance)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

func (h DBHandler) FindBalanceHistoryByUser(username string) (*[]models.BalanceHistory, int) {
	var balance []models.BalanceHistory
	user, err := h.GetUserByUsername(username)
	if err != -1 {
		return nil, err
	}
	res := h.DB.Model(&user).Association("BalanceHistory").Find(&balance)
	if res != nil {
		return nil, dbHandleError(res)
	}
	return &balance, -1
}

func (h DBHandler) AddBalanceHistory(balance models.BalanceHistory, userId string) int {
	user, err := h.GetUserByUsername(userId)
	if err != -1 {
		return err
	}
	res := h.DB.Model(&user).Association("BalanceHistory").Append(&balance)
	if res != nil {
		return dbHandleError(res)
	}
	return -1
}

// Bet methods

func (h DBHandler) CreateBet(bet models.Bet) int {
	res := h.DB.Create(&bet)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

func (h DBHandler) FindBet(betID int) (*models.Bet, int) {
	var bet models.Bet
	res := h.DB.First(&bet, betID)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &bet, -1
}

func (h DBHandler) UpdateBet(bet models.Bet) int {
	res := h.DB.Save(&bet)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

func (h DBHandler) DeleteBet(betID int) int {
	res := h.DB.Delete(&models.Bet{}, betID)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

func (h DBHandler) GetUserBet(username string) (*[]models.UserBet, int) {
	var user models.User
	res := h.DB.Where("username = ?", username).First(&user)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}

	var bets []models.UserBet
	res = h.DB.Where("user_id = ?", user.ID).Find(&bets)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}

	return &bets, -1
}

func (h DBHandler) GetBetsByBetID(betID uint) (*[]models.UserBet, int) {
	var bets []models.UserBet
	res := h.DB.Where("bet_id = ?", betID).Find(&bets)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &bets, -1
}

// Helper functions

func dbHandleError(e error) int {
	var res int
	if errors.Is(e, gorm.ErrDuplicatedKey) {
		res = tools.DB_DUP_KEY
	} else if errors.Is(e, gorm.ErrRecordNotFound) {
		res = tools.DB_REC_NOTFOUND
	} else {
		res = tools.DB_UNKOWN_ERR
	}
	log.Info("DB Error: ", e.Error(), res)
	return res
}
