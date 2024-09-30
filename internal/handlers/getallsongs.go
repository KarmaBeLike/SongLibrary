package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/SongLibrary/internal/service"
)

type SongHandler struct {
	songService *service.SongService
}

func NewSongHandler(songService *service.SongService) *SongHandler {
	return &SongHandler{songService: songService}
}

func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	songs, pagination, err := h.songService.GetSongs(group, song, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"songs":      songs,
		"pagination": pagination,
	})
}
