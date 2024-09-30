package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid song ID. It must be a positive integer.", http.StatusBadRequest)
		return
	}

	var updateRequest models.UpdateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	if err := h.songService.UpdateSongByID(id, &updateRequest); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "Song updated successfully.",
		"id":      id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
