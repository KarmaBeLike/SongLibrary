package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

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
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received request to add a new song", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	var newSong models.NewSongRequest

	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		slog.Error("Invalid request payload", slog.Any("error", err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	slog.Debug("Request payload decoded", slog.Any("newSong", newSong))

	id, err := h.songService.AddSong(r.Context(), newSong)
	if err != nil {
		slog.Error("Failed to add song", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Song added successfully", slog.Int("songID", id))

	response := map[string]interface{}{
		"message": "Song added successfully",
		"id":      id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Failed to encode response", slog.Any("error", err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	slog.Debug("Response sent", slog.Int("songID", id))
}
