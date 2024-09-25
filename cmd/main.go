package main

import (
	"log"
	"net/http"

	"github.com/KarmaBeLike/SongLibrary/internal/routers"
)

func main() {
	router := routers.SetupRoutes()

	log.Println("Server is runnig on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}
