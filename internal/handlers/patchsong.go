package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
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

	if err := h.songService.UpdateSongByID(id, &updateRequest); err != nil {
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
