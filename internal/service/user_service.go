package service

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
)

func CheckUserNamePassword(username string, pass string) bool {
	db := database.NewDatabaseConnection()
	var user uint
	db.Select("id").Where("username = ? and password_hash = ", username, pass).Find(&user)
	return user == 0
}

func SaveUser(userDto *common.UserReg) model.User {

	// create user
	user := model.User{
		Username:     userDto.Username,
		Email:        userDto.Email,
		PasswordHash: userDto.Password,
	}

	db := database.NewDatabaseConnection()
	db.Save(&user)

	return user
}
