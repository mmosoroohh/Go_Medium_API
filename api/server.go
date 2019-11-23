package api

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mmosoroohh/Go_Medium_API/api/controllers"
	"github.com/mmosoroohh/Go_Medium_API/api/seed"
	"log"
)

var server = controllers.Server{}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error connecting env, %v", err)
	} else {
		fmt.Println("Now connecting to the env")
	}

	server.Initialize()

	seed.Load(server.DB)

	server.Run(":8080")
}
