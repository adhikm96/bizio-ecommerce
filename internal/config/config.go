package config

import (
	"os"
	"strconv"
)

var Config = map[string]string{
	"DB_NAME":     os.Getenv("DB_NAME"),
	"DB_USERNAME": os.Getenv("DB_USERNAME"),
	"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
	"DB_PORT":     os.Getenv("DB_PORT"),
	"DB_HOST":     os.Getenv("DB_HOST"),
}

func Get(key string) string {
	return Config[key]
}

func GetInt(key string) int {
	data, _ := strconv.Atoi(Config[key])
	return data
}
