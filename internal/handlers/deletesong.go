package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		slog.Warn("Invalid song ID received", slog.String("id", idStr))
		http.Error(w, "Invalid song ID. It must be a positive integer.", http.StatusBadRequest)
		return
	}

	slog.Info("Starting to delete song", slog.Int("song_id", id))

	err = h.songService.DeleteSongByID(id)
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
