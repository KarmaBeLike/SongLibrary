package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *SongHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received request to get song by ID", slog.String("method", r.Method), slog.String("url", r.URL.String()))

	idStr := r.URL.Query().Get("id")
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

	response, err := h.songService.GetPaginatedSongLyrics(id, page, limit)
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
