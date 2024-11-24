package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
	"github.com/KarmaBeLike/SongLibrary/internal/service"
)

type SongHandler struct {
	songService *service.SongService
}

func NewSongHandler(songService *service.SongService) *SongHandler {
	return &SongHandler{songService: songService}
}

// GetSongs fetches songs with optional filters (group and song title) and supports pagination.
// @Summary Get songs with optional filtering and pagination
// @Description Retrieve a list of songs filtered by group or title, with pagination support.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Group name" example("The Beatles")
// @Param song query string false "Song title" example("Hey Jude")
// @Param page query int false "Page number" example(1)
// @Param limit query int false "Limit of songs per page" example(10)
// @Success 200 {object} map[string]interface{} "Successful operation"
// @Failure 500 {string} string "Internal server error"
// @Router /api/songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	// Создаём экземпляр структуры фильтров
	filter := &models.SongFilter{
		Group:       r.URL.Query().Get("group"),
		Song:       r.URL.Query().Get("song"),
		Text:        r.URL.Query().Get("text"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
	}

	// Обрабатываем параметры `page` и `limit`
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
		slog.Debug("Invalid page number, setting to default", slog.Int("page", page))
	}
	filter.Page = page

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
		slog.Debug("Invalid limit, setting to default", slog.Int("limit", limit))
	}
	filter.Limit = limit

	slog.Debug("Received filter request", slog.Any("filter", filter))
	// Вызов сервиса с фильтром

	songs, pagination, err := h.songService.GetSongs(r.Context(), *filter)
	if err != nil {
		slog.Error("Failed to fetch songs", slog.Any("filter", filter), slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Songs fetched successfully", slog.Int("song_count", len(songs)), slog.Any("pagination", pagination))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"songs":      songs,
		"pagination": pagination,
	})
}
