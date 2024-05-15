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
	defaultDbName     = "billbook"
	defaultDbPort     = "5432"
	defaultDbHost     = "localhost"
	defaultDbPassword = "password"
	defaultServerHost = "8080"
	defaultDbUser     = "seller"
	defaultRate		  = 5
)

var (
	dbName     string
	dbPort     string
	dbHost     string
	dbPassword string
	dbUser     string
	serverHost string
	rate 	   int
)

func main(){
	dbname:= os.Getenv("DB_NAME")
	if len(dbname)==0 {
		dbname=defaultDbName
	}
	dbName=dbname
	dbport:= os.Getenv("DB_PORT")
	if len(dbport)==0 {
		dbport=defaultDbPort
	}
	dbPort=dbport

	dbhost:= os.Getenv("DB_HOST")
	if len(dbhost)==0 {
		dbhost=defaultDbHost
	}
	dbHost=dbhost

	dbuser:= os.Getenv("DB_USER")
	if len(dbuser)==0 {
		dbuser=defaultDbUser
	}
	dbUser=dbuser

	dbpassword:= os.Getenv("DB_PASSWORD")
	if len(dbpassword)==0 {
		dbpassword=defaultDbPassword
	}
	dbPassword=dbpassword

	serverhost:=os.Getenv("SERVER_HOST")
	if len(serverhost)==0 {
		serverhost=defaultServerHost
	}
	serverHost=serverhost

	limitrate,_:=strconv.Atoi(os.Getenv("RATE"))
	if limitrate==0 {
		limitrate=defaultRate
	}
	rate=limitrate
	
	db.Connect(dbHost, dbPort, dbUser, dbName, dbPassword)

	if err := db.DB.AutoMigrate(&model.Order{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	r := apis.Router()

	if err := r.Run(); err != nil {
		log.Fatalf("Error while running the server: %v", err)
	}
}