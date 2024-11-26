package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/KarmaBeLike/SongLibrary/config"
	_ "github.com/KarmaBeLike/SongLibrary/docs"
	"github.com/KarmaBeLike/SongLibrary/internal/api"
	"github.com/KarmaBeLike/SongLibrary/internal/database"
	"github.com/KarmaBeLike/SongLibrary/internal/repository"
	"github.com/KarmaBeLike/SongLibrary/internal/routers"
	"github.com/KarmaBeLike/SongLibrary/internal/service"
)

// @title SongLibrary
// @version 1.0

// @host localhost:8080
// @BasePath /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	slog.Info("Logger initialized", slog.String("output", "JSON"))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return
	}

	db, err := database.OpenDB(cfg)
	if err != nil {
		slog.Error("failed to connect to db", slog.Any("error", err))
		return
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		slog.Error("error running migrations", slog.Any("error", err))
		return
	}
	err = database.LoadTestData(db, "migrations/seed_data.sql")
	if err != nil {
		log.Fatal("Error loading test data:", err)
	}

	songRepo := repository.NewSongRepository(db)
	apiClient := api.NewExternalAPI(cfg.ExternalAPI)
	fmt.Println(cfg.ExternalAPI)
	songService := service.NewSongService(songRepo, apiClient)
	router := routers.SetupRoutes(songService)

	port := cfg.Port
	log.Printf("Server is running on port %d...", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
