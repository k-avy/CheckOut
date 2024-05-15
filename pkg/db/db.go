package db

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RED *redis.Client
)

func ConnectPG(dbhost, port, user, dbname, password string) {

	dburi := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbhost, port, user, dbname, password)
	gormConfig := &gorm.Config{}
	database, err := gorm.Open(postgres.Open(dburi), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = database
	log.Println("Postgres DB connected")
}

func ConnectRedis(rhost string, rport string) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", rhost, rport),
		DB:   0,
	})
	RED = client
	log.Println("Redis DB connected")
}
