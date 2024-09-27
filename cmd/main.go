package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/KarmaBeLike/SongLibrary/config"
	"github.com/KarmaBeLike/SongLibrary/internal/database"
	"github.com/KarmaBeLike/SongLibrary/internal/routers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config loh", slog.Any("error", err))
		return
	}

	db, err := database.OpenDB(cfg)
	if err != nil {
		slog.Error("bd loh", slog.Any("error", err))
		return
	}
	defer db.Close()
	slog.Info("health check", slog.Any("db ping", db.Ping()))

	router := routers.SetupRoutes()

	log.Println("Server is runnig on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}
