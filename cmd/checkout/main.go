package main

import "os"

const (
	defaultDbName     = "billbook"
	defaultDbPort     = "5432"
	defaultDbHost     = "localhost"
	defaultDbPassword = "password"
	defaultServerHost = "8080"
	defaultDbUser     = "seller"

)

var (
	dbName     string
	dbPort     string
	dbHost     string
	dbPassword string
	dbUser     string
	serverHost string
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
}