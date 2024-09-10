package handlers

import (
	"errors"
	"fmt"
	"gambler/backend/database"
	"gambler/backend/database/models"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/tools"
	"runtime"

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
		return dbHandleError(errors.New("user ID not populated after creation"))
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

// BalanceHistory methods

func (h DBHandler) UpdateUserBalance(amount float64, user models.User, reason string) int {
	user.Balance += amount
	res := h.DB.Save(&user)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	err := h.AddBalanceHistory(models.BalanceHistory{
		UserID: user.ID,
		Amount: amount,
		Reason: reason,
	}, user.ID)
	if err != -1 {
		return err
	}
	return -1
}

func (h DBHandler) CreateBalanceHistory(balance models.BalanceHistory) int {
	res := h.DB.Create(&balance)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

func (h DBHandler) FindBalanceHistoryByUser(userId uint) (*[]models.BalanceHistory, int) {
	var balance []models.BalanceHistory
	user, err := h.GetUserByID(userId)
	if err != -1 {
		return nil, err
	}
	res := h.DB.Model(&user).Association("BalanceHistory").Find(&balance)
	if res != nil {
		return nil, dbHandleError(res)
	}
	return &balance, -1
}

func (h DBHandler) AddBalanceHistory(balance models.BalanceHistory, userId uint) int {
	user, err := h.GetUserByID(userId)
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

func (h DBHandler) CreateBet(bet models.Bet, userId uint, betOption string, amount float64) int {
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the user
	if err := tx.Create(&bet).Error; err != nil {
		tx.Rollback()
		log.Errorf("Error creating user: %v", err)
		return dbHandleError(err)
	}

	// Get the user ID
	user, err := h.GetUserByID(userId)
	if err != -1 {
		tx.Rollback()
		return err
	}

	// Ensure the user ID is populated
	if bet.ID == 0 {
		tx.Rollback()
		log.Error("Bet ID not populated after creation")
		return dbHandleError(errors.New("bet ID not populated after creation"))
	}

	// Create the initial balance history entry
	initialBet := models.UserBet{
		UserID:    user.ID,
		BetID:     bet.ID,
		Amount:    amount,
		BetOption: betOption,
	}

	if err := tx.Create(&initialBet).Error; err != nil {
		tx.Rollback()
		log.Errorf("Error creating balance history: %v", err)
		return dbHandleError(err)
	}

	user.Balance -= amount
	res := tx.Save(&user)
	if res.Error != nil {
		tx.Rollback()
		log.Errorf("Error updating user balance: %v", res.Error)
		return dbHandleError(res.Error)
	}

	balance := models.BalanceHistory{
		UserID: user.ID,
		Amount: amount,
		Reason: fmt.Sprintf("Bet on: %s", bet.Name),
	}

	bRes := tx.Model(&user).Association("BalanceHistory").Append(&balance)
	if bRes != nil {
		return dbHandleError(bRes)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return dbHandleError(err)
	}

	err = Cache.UpdateBet(bet.ID)
	if err != -1 {
		return err
	}

	return -1
}

func (h DBHandler) FindBet(betID int) (*models.Bet, int) {
	var bet models.Bet
	res := h.DB.Preload("UserBets").Where("deleted_at IS NULL").First(&bet, betID)
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

func (h DBHandler) GetUserBet(userId uint) (*[]models.UserBet, int) {
	var user models.User
	res := h.DB.Where("ID = ?", userId).First(&user)
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

func (h DBHandler) GetUserBetByID(id uint) (*models.UserBet, int) {
	log.Info(id)
	var bet models.UserBet
	res := h.DB.First(&bet, id)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &bet, -1
}

func (h DBHandler) GetUserBetByBetID(betID uint, userId uint) (*models.UserBet, int) {
	var bet models.UserBet
	res := h.DB.Preload("UserBets").Where("user_id = ? AND bet_id = ? AND deleted_at IS NULL", userId, betID).First(&bet)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}

	return &bet, -1
}

func (h DBHandler) GetBetsByBetID(betID uint) (*[]models.UserBet, int) {
	var bets []models.UserBet
	res := h.DB.Preload("UserBets").Where("bet_id = ? AND deleted_at IS NULL", betID).Find(&bets)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &bets, -1
}

func (h DBHandler) GetBetByID(id uint) (*models.Bet, int) {
	var bet models.Bet
	res := h.DB.Preload("UserBets").First(&bet, id)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &bet, -1
}

func (h DBHandler) GetBetByBetName(name string) (*models.Bet, int) {
	var bet models.Bet
	res := h.DB.Preload("UserBets").Where("name = ?", name).First(&bet)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &bet, -1
}

func (h DBHandler) PlaceBet(userBet models.UserBet) int {
	res := h.DB.Create(&userBet)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}
	return -1
}

func (h DBHandler) CancelBet(userBet models.UserBet, user models.User) int {
	res := h.DB.Delete(&models.UserBet{}, userBet.ID)
	if res.Error != nil {
		return dbHandleError(res.Error)
	}

	return -1
}

func (h DBHandler) GetAllBets() (*[]models.Bet, int) {
	var bet []models.Bet
	res := h.DB.Preload("UserBets").Find(&bet)
	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}
	return &bet, -1
}

func (h DBHandler) GetAllActiveBets() (*[]models.Bet, int) {
	var bets []models.Bet

	// Use Preload to also load associated UserBets for each Bet
	res := h.DB.Where("status = ?", customTypes.Open).
		Preload("UserBets"). // Preload the UserBets relation
		Find(&bets)

	if res.Error != nil {
		log.Info(res.Error)
		return nil, dbHandleError(res.Error)
	}

	return &bets, -1
}

func (h DBHandler) GetAllClosedBets() (*[]models.Bet, int) {
	var bets []models.Bet

	// Use Preload to also load associated UserBets for each Bet
	res := h.DB.Where("status = ?", customTypes.Closed).
		Preload("UserBets"). // Preload the UserBets relation
		Find(&bets)

	if res.Error != nil {
		return nil, dbHandleError(res.Error)
	}

	return &bets, -1
}

// Helper functions

func dbHandleError(e error) int {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Info(fmt.Sprintf("Called from %s, line %d", file, line))
	}
	var res int
	if errors.Is(e, gorm.ErrDuplicatedKey) {
		res = tools.DB_DUP_KEY
	} else if errors.Is(e, gorm.ErrRecordNotFound) {
		res = tools.DB_REC_NOTFOUND
	} else {
		res = tools.DB_UNKNOWN_ERR
	}
	log.Info("DB Error: ", e.Error(), res)
	return res
}
