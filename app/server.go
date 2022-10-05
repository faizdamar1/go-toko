package app

import (
	"flag"
	"log"
	"os"

	"github.com/faizdamar1/go-toko/app/controllers"
	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Run() {
	var server = controllers.Server{}
	var appConfig = controllers.AppConfig{}
	var dbConfig = controllers.DBConfig{}

	err := godotenv.Load()

	if err != nil {
		log.Fatalf(err.Error())
	}

	appConfig.AppName = getEnv("APP_NAME", "GoTokoApp")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	dbConfig.DBDriver = getEnv("DB_HOST", "localhost")
	dbConfig.DBName = getEnv("DB_NAME", "dbname")
	dbConfig.DBUser = getEnv("DB_USER", "dbuser")
	dbConfig.DBPassword = getEnv("DB_PASS", "")
	dbConfig.DBPort = getEnv("DB_PORT", "3306")
	dbConfig.DBDriver = getEnv("DB_DRIVER", "mysql")

	flag.Parse()
	arg := flag.Arg(0)

	if arg != "" {
		server.InitCommands(appConfig, dbConfig)
	} else {
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}

}
