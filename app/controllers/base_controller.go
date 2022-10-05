package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/faizdamar1/go-toko/app/models"
	"github.com/faizdamar1/go-toko/database/seeders"
	"github.com/gorilla/mux"
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

func (server *Server) dbMigrate() {
	for _, model := range models.RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Migrate successfull")
}

func (server *Server) InitCommands(config AppConfig, dbConfig DBConfig) {
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

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("welcome to " + appConfig.AppName)
	server.initializeRoutes()
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

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
