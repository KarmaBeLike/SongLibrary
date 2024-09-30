package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var newSong models.NewSongRequest

	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := h.songService.AddSong(newSong)
	if err != nil {
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
