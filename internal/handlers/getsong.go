package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

func GetSongHandler(w http.ResponseWriter, r *http.Request) {
	// Пример текста песни
	lyrics := `Ooh baby, don't you know I suffer?
Ooh baby, can you hear me moan?
You caught me under false pretenses
How long before you let me go?

Ooh
You set my soul alight
Ooh
You set my soul alight

Glaciers melting in the dead of night (ooh)
And the superstars sucked into the supermassive (you set my soul alight)
Glaciers melting in the dead of night
And the superstars sucked into the (you set my soul)
(Into the supermassive)

I thought I was a fool for no one
Ooh baby, I'm a fool for you
You're the queen of the superficial
And how long before you tell the truth?

Ooh
You set my soul alight
Ooh
You set my soul alight

Glaciers melting in the dead of night (ooh)
And the superstars sucked into the supermassive (you set my soul alight)
Glaciers melting in the dead of night
And the superstars sucked into the (you set my soul)
(Into the supermassive)

Supermassive black hole
Supermassive black hole
Supermassive black hole
Supermassive black hole
Glaciers melting in the dead of night
And the superstars sucked into the supermassive
Glaciers melting in the dead of night
And the superstars sucked into the supermassive
Glaciers melting in the dead of night (ooh)
And the superstars sucked into the supermassive (you set my soul alight)
Glaciers melting in the dead of night
And the superstars sucked into the (you set my soul)
(Into the supermassive)

Supermassive black hole
Supermassive black hole
Supermassive black hole
Supermassive black hole`

	// Разделяем текст песни на куплеты по двум переносам строки
	verses := strings.Split(lyrics, "\n\n")

	for i := range verses {
		verses[i] = strings.ReplaceAll(verses[i], "\n", " ")
	}

	// Получаем параметры запроса: page и limit (если они есть)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Значения по умолчанию для пагинации
	page := 1
	limit := 1 // сколько куплетов возвращать на страницу по умолчанию

	// Если параметр страницы передан, конвертируем его в int
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p <= 0 {
			http.Error(w, "Invalid page number. It must be a positive integer.", http.StatusBadRequest)
			return
		}
		page = p
	}

	// Если параметр лимита передан, конвертируем его в int
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			http.Error(w, "Invalid limit value. It must be a positive integer.", http.StatusBadRequest)
			return
		}
		limit = l
	}

	// Вычисляем сколько всего куплетов
	totalVerses := len(verses)

	// Вычисляем начальный и конечный индексы для пагинации
	start := (page - 1) * limit
	end := start + limit

	if start >= totalVerses || start < 0 {
		http.Error(w, "Page out of range.", http.StatusBadRequest)
		return
	}

	// Извлекаем куплеты в соответствии с номером страницы

	if end > totalVerses {
		end = totalVerses
	}
	paginatedVerses := verses[start:end]

	// Формируем структуру ответа
	song := models.Song{
		Group:  "Muse",
		ID:     1,
		Song:   "Supermassive Black Hole",
		Limit:  limit,
		Page:   page,
		Total:  totalVerses,
		Verses: paginatedVerses,
	}

	// Возвращаем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}
