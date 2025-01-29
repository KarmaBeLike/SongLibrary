package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

type (
	songService interface {
		GetSongs(ctx context.Context, filter models.SongFilter) ([]models.Song, models.Pagination, error)
		GetPaginatedSongLyrics(ctx context.Context, id int, page, limit int) (*models.SongVerses, error)
		DeleteSongByID(ctx context.Context, id int) error
		UpdateSongByID(ctx context.Context, id int, updateRequest *models.UpdateSongRequest) error
		AddSong(ctx context.Context, newSong models.NewSongRequest) (int, error)
	}
	SongClient struct {
		service songService
	}
)

func NewSongClient(service songService) *SongClient {
	return &SongClient{
		service: service,
	}
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
func (c *SongClient) GetSongs(w http.ResponseWriter, r *http.Request) {
	// Создаём экземпляр структуры фильтров
	filter := &models.SongFilter{
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		Text:        r.URL.Query().Get("text"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Link:        r.URL.Query().Get("link"),
	}

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

	songs, pagination, err := c.service.GetSongs(r.Context(), *filter)
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

// GetSongByID retrieves the lyrics of a song by its ID with optional pagination.
// @Summary Get song lyrics by ID with optional pagination
// @Description Retrieve lyrics of a song by its ID, with options to paginate the lyrics.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id query int true "Song ID" example(1)
// @Param page query int false "Page number" example(1)
// @Param limit query int false "Limit of verses per page" example(1)
// @Success 200 {object} models.SongVerses "Successful operation"
// @Failure 400 {string} string "Invalid song ID or pagination parameters"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/songs/lyrics [get]
func (c *SongClient) GetSongLyrics(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received request to get song by ID", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	idStr := r.URL.Query().Get("song_id")

	if idStr == "" {
		slog.Error("Invalid song ID", slog.Any("idStr", idStr))
		http.Error(w, "Missing or invalid song_id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		slog.Error("Invalid song ID", slog.String("idStr", idStr), slog.Any("error", err))
		http.Error(w, "Invalid song ID.", http.StatusBadRequest)
		return
	}
	slog.Info("Fetching song by ID", slog.Int("id", id))

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 1

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)

		if err != nil || page <= 0 {
			slog.Error("Invalid page number", slog.String("pageStr", pageStr), slog.Any("error", err))
			http.Error(w, "Invalid page number. It must be a positive integer.", http.StatusBadRequest)
			return
		}
		slog.Debug("Page number set", slog.Int("page", page))
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)

		if err != nil || limit <= 0 {
			slog.Error("Invalid limit value", slog.String("limitStr", limitStr), slog.Any("error", err))
			http.Error(w, "Invalid limit value. It must be a positive integer.", http.StatusBadRequest)
			return
		}
		slog.Debug("Limit set", slog.Int("limit", limit))
	}

	slog.Info("Fetching paginated song lyrics", slog.Int("id", id), slog.Int("page", page), slog.Int("limit", limit))

	response, err := c.service.GetPaginatedSongLyrics(r.Context(), id, page, limit)
	if err != nil {
		slog.Error("Failed to fetch song lyrics", slog.Int("id", id), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	slog.Info("Song fetched successfully", slog.Int("id", id))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Failed to encode response", slog.Any("error", err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	slog.Debug("Response sent", slog.Int("id", id), slog.Int("page", page), slog.Int("limit", limit))
}

// AddSong adds a new song to the library.
// @Summary Add a new song
// @Description Adds a new song with the given details to the library.
// @Tags songs
// @Accept json
// @Produce json
// @Param newSong body models.NewSongRequest true "New song details"
// @Success 201 {object} map[string]interface{} "Successful operation"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /api/songs [post]
func (c *SongClient) AddSong(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received request to add a new song", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	var newSong models.NewSongRequest

	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		slog.Error("Invalid request payload", slog.Any("error", err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := c.service.AddSong(r.Context(), newSong)
	if err != nil {
		slog.Error("Failed to add song", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Song added successfully",
		"id":      id,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateSong updates the details of a song by its ID.
// @Summary Update song by ID
// @Description Update the details of a song using its ID.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id query int true "Song ID" example(1)
// @Param updateRequest body models.UpdateSongRequest true "Update song details"
// @Success 200 {object} map[string]interface{} "Successful operation"
// @Failure 400 {string} string "Invalid song ID or request body"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/songs [patch]
func (c *SongClient) UpdateSong(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received request to update song", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		slog.Error("Invalid song ID", slog.String("idStr", idStr), slog.Any("error", err))
		http.Error(w, "Invalid song ID. It must be a positive integer.", http.StatusBadRequest)
		return
	}
	slog.Info("Updating song by ID", slog.Int("id", id))

	var updateRequest models.UpdateSongRequest

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		slog.Error("Failed to decode request body", slog.Any("error", err))
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}
	slog.Debug("Request body decoded", slog.Any("updateRequest", updateRequest))

	if err := c.service.UpdateSongByID(r.Context(), id, &updateRequest); err != nil {
		slog.Error("Failed to update song", slog.Int("id", id), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Song updated successfully", slog.Int("id", id))

	response := map[string]interface{}{
		"message": "Song updated successfully.",
		"id":      id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Failed to encode response", slog.Any("error", err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	slog.Debug("Response sent", slog.Int("id", id))
}

// DeleteSong deletes a song by its ID
// @Summary Delete a song by ID
// @Description Delete a song from the library by providing the song ID as a query parameter.
// @Tags songs
// @Accept json
// @Produce json
// @Param id query int true "Song ID"
// @Success 200 {object} map[string]interface{} "Song deleted successfully"
// @Failure 400 {string} string "Invalid song ID. It must be a positive integer."
// @Failure 404 {string} string "Song not found"
// @Router /api/songs [delete]
func (c *SongClient) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		slog.Warn("Invalid song ID received", slog.String("id", idStr))
		http.Error(w, "Invalid song ID. It must be a positive integer.", http.StatusBadRequest)
		return
	}

	slog.Info("Starting to delete song", slog.Int("song_id", id))

	err = c.service.DeleteSongByID(r.Context(), id)
	if err != nil {
		slog.Error("Failed to delete song", slog.Int("song_id", id), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Song deleted successfully", slog.Int("song_id", id))

	response := map[string]interface{}{
		"message": "Song deleted successfully.",
		"id":      id,
	}

	slog.Debug("Sending delete song response", slog.Any("response", response))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
