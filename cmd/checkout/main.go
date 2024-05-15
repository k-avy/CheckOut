package main

import (
	"log"
	"os"
	"strconv"

	"github.com/k-avy/CheckOut/pkg/apis"
	"github.com/k-avy/CheckOut/pkg/db"
	"github.com/k-avy/CheckOut/pkg/model"
)

const (
	defaultDbName     = "postgres"
	defaultDbPort     = "5432"
	defaultDbHost     = "localhost"
	defaultDbPassword = "password"
	defaultServerPort = "8080"
	defaultDbUser     = "seller"
	defaultRedishost  = "localhost"
	defaultRedisPort  = "6379"
	defaultRate       = 5
)

var (
	dbName     string
	dbPort     string
	dbHost     string
	dbPassword string
	dbUser     string
	serverPort string
	redisHost  string
	redisPort  string
	rate       int64
)

func main() {
	dbname := os.Getenv("DB_NAME")
	if len(dbname) == 0 {
		dbname = defaultDbName
	}
	dbName = dbname

	dbport := os.Getenv("DB_PORT")
	if len(dbport) == 0 {
		dbport = defaultDbPort
	}
	dbPort = dbport

	dbhost := os.Getenv("DB_HOST")
	if len(dbhost) == 0 {
		dbhost = defaultDbHost
	}
	dbHost = dbhost

	dbuser := os.Getenv("DB_USER")
	if len(dbuser) == 0 {
		dbuser = defaultDbUser
	}
	dbUser = dbuser

	dbpassword := os.Getenv("DB_PASSWORD")
	if len(dbpassword) == 0 {
		dbpassword = defaultDbPassword
	}
	dbPassword = dbpassword

	serverport := os.Getenv("SERVER_PORT")
	if len(serverport) == 0 {
		serverport = defaultServerPort
	}
	serverPort = serverport

	limitrate, _ := strconv.Atoi(os.Getenv("RATE"))
	if limitrate == 0 {
		limitrate = defaultRate
	}
	rate = int64(limitrate)

	redisport := os.Getenv("REDIS_PORT")
	if len(redisport) == 0 {
		redisport = defaultRedisPort
	}
	redisPort = redisport

	redishost := os.Getenv("REDIS_HOST")
	if len(redishost) == 0 {
		redishost = defaultRedishost
	}
	redisHost = redishost

	db.ConnectPG(dbHost, dbPort, dbUser, dbName, dbPassword)
	db.ConnectRedis(redisHost, redisPort)

	if err := db.DB.AutoMigrate(&model.Order{}, &model.User{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	r := apis.Router(rate)

	if err := r.Run(":" + serverPort); err != nil {
		log.Fatalf("Error while running the server: %v", err)
	}
}
