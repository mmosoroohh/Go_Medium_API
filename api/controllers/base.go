package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/mmosoroohh/Go_Medium_API/api/models"
	"log"
	"net/http"
	"os"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDriver := os.Getenv("DB_DRIVER")

	if dbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
		server.DB, err = gorm.Open(dbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", dbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are now connected to %s database", dbDriver)
		}
	}

	if dbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)
		server.DB, err = gorm.Open(dbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", dbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are now connected to the %s database", dbDriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) // Database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
