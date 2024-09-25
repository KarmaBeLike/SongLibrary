package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KarmaBeLike/SongLibrary/internal/services"
)

func SongsHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")    // Фильтрация по группе
	song := r.URL.Query().Get("song")      // Фильтрация по названию песни
	pageStr := r.URL.Query().Get("page")   // Номер страницы
	limitStr := r.URL.Query().Get("limit") // Количество элементов на странице

	// Параметры по умолчанию для пагинации
	page := 1   // Страница по умолчанию
	limit := 10 // Лимит по умолчанию (10 песен на страницу)

	// Если параметр страницы передан, конвертируем его в int
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	// Если параметр лимита передан, конвертируем его в int
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	// Вызов сервисного слоя для получения песен с фильтрацией и пагинацией
	songs, total, err := services.GetSongs(group, song, page, limit)
	if err != nil {
		http.Error(w, "Error fetching songs", http.StatusInternalServerError)
		return
	}

	// Формирование ответа с данными о песнях и пагинации
	response := map[string]interface{}{
		"songs": songs,
		"page":  page,
		"limit": limit,
		"total": total, // Общее количество песен, доступное для отображения
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
