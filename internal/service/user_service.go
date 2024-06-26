package service

import "github.com/Digital-AIR/bizio-ecommerce/internal/database"

func CheckUserExists(uId uint) bool {
	db := database.NewDatabaseConnection()

	var userExists bool
	db.Raw("select exists(select 1 from users where id = ?)", uId).Scan(&userExists)

	return userExists
}
