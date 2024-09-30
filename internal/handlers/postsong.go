package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received request to add a new song", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	var newSong models.NewSongRequest

	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		slog.Error("Invalid request payload", slog.Any("error", err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	slog.Debug("Request payload decoded", slog.Any("newSong", newSong))

	id, err := h.songService.AddSong(newSong)
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
