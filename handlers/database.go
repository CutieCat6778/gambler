package handlers

import (
	"errors"
	"gambler/backend/database"
	"gambler/backend/database/models"
	"gambler/backend/tools"

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

func (h DBHandler) CreateUser(user models.User) int {
	res := h.DB.Create(&user)
	if res.Error != nil {
		return HandleError(res.Error)
	}
	return -1
}

func (h DBHandler) UpdateUser(user models.User) (*models.User, int) {
	res := h.DB.Save(&user)
	if res.Error != nil {
		return nil, HandleError(res.Error)
	}
	return &user, -1
}

func (h DBHandler) GetUserByID(id uint) (*models.User, int) {
	var user models.User
	res := h.DB.First(&user, id)
	if res.Error != nil {
		return nil, HandleError(res.Error)
	}
	return &user, -1
}

func (h DBHandler) GetUserByUsername(username string) (*models.User, int) {
	var user models.User
	res := h.DB.Where("username = ?", username).First(&user)
	if res.Error != nil {
		return nil, HandleError(res.Error)
	}
	return &user, -1
}

func (h DBHandler) DeleteUserByID(id uint) int {
	res := h.DB.Delete(&models.User{}, id)
	if res.Error != nil {
		return HandleError(res.Error)
	}
	return -1
}

func HandleError(e error) int {
	if errors.Is(e, gorm.ErrDuplicatedKey) {
		return tools.DB_DUP_KEY
	} else if errors.Is(e, gorm.ErrRecordNotFound) {
		return tools.DB_REC_NOTFOUND
	} else {
		return tools.DB_UNKOWN_ERR
	}
}
