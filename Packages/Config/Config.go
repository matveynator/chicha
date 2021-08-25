package Config

import (
	"fmt"
	"github.com/joho/godotenv" // Enviroment read package
	"log"
	"os"
	"strconv"
)

var APP_ANTENNA_LISTENER_IP,API_SERVER_LISTENER_IP,TIME_ZONE,RTS,DB_TYPE,DB_HOST,DB_USER,DB_PASSWORD,DB_NAME,DB_PORT,ADMIN_LOGIN,ADMIN_PASSWORD,MINIMAL_LAP_TIME,LAPS_SAVE_INTERVAL,PROXY_ACTIVE,PROXY_HOST,PROXY_PORT string
var RTS8 int

func init()  {
	// Init ConfigMap here
	// Load enviroment
	// PROXY settings
	fmt.Println("Load enviroment")
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}

	// PROXY settings
	PROXY_ACTIVE, _ = os.LookupEnv("PROXY_ACTIVE")
	PROXY_HOST, _ = os.LookupEnv("PROXY_HOST")
	PROXY_PORT, _ = os.LookupEnv("PROXY_PORT")
	// Check enviroment
	APP_ANTENNA_LISTENER_IP, _ = os.LookupEnv("APP_ANTENNA_LISTENER_IP")
	API_SERVER_LISTENER_IP, _ = os.LookupEnv("API_SERVER_LISTENER_IP")
	TIME_ZONE, _ =  os.LookupEnv("TIME_ZONE")
	RTS, _ = os.LookupEnv("RACE_TIMEOUT_SEC")
	RTS8, _ = strconv.Atoi(RTS)

	// DB connection preferences
	DB_HOST, _ = os.LookupEnv("DB_HOST")
	DB_USER, _ = os.LookupEnv("DB_USER")
	DB_PASSWORD, _ = os.LookupEnv("DB_PASSWORD")
	DB_NAME, _ = os.LookupEnv("DB_NAME")
	DB_PORT, _ = os.LookupEnv("DB_PORT")
	DB_TYPE, _ = os.LookupEnv("DB_TYPE")

	ADMIN_LOGIN, _ = os.LookupEnv("ADMIN_LOGIN")
	ADMIN_PASSWORD, _ = os.LookupEnv("ADMIN_PASSWORD")

	MINIMAL_LAP_TIME, _ = os.LookupEnv("MINIMAL_LAP_TIME")
	LAPS_SAVE_INTERVAL, _ = os.LookupEnv("LAPS_SAVE_INTERVAL")

	return
}
