package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/faizdamar1/go-toko/database/seeders"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

func (server *Server) initCommands(config AppConfig, dbConfig DBConfig) {
	server.InitializeDB(dbConfig)

	cmdApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "db:migrate",
				Aliases: []string{"a"},
				Usage:   "make a migration database",
				Action: func(cCtx *cli.Context) error {
					server.dbMigrate()
					return nil
				},
			},
			{
				Name:    "db:seed",
				Aliases: []string{"a"},
				Usage:   "make a seed to database",
				Action: func(cCtx *cli.Context) error {
					err := seeders.DBSeed(server.DB)

					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
		}}

	if err := cmdApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func (server *Server) Inisialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("welcome to " + appConfig.AppName)
	server.inisializeRoutes()
}

func (server *Server) InitializeDB(dbConfig DBConfig) {
	var err error
	if dbConfig.DBDriver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Panic(err.Error())
	} else {
		fmt.Println("Connected to DB")
	}
}

func (server *Server) dbMigrate() {
	for _, model := range RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Migrate successfull")
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

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
		server.initCommands(appConfig, dbConfig)
	} else {
		server.Inisialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}

}
