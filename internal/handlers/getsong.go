package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *SongHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid song ID. It must be a positive integer.", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			http.Error(w, "Invalid page number. It must be a positive integer.", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit value. It must be a positive integer.", http.StatusBadRequest)
			return
		}
	}

	response, err := h.songService.GetPaginatedSongLyrics(id, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
