package services

import "github.com/KarmaBeLike/SongLibrary/internal/models"

func GetSongs(group, song string, page, limit int) ([]models.Song, int, error) {
	// Это пример заглушки, на реальном проекте здесь была бы работа с базой данных
	allSongs := []models.Song{
		{Group: "Muse", Song: "Supermassive Black Hole"},
		{Group: "Radiohead", Song: "Creep"},
		{Group: "Nirvana", Song: "Smells Like Teen Spirit"},
		{Group: "The Beatles", Song: "Hey Jude"},
		{Group: "Muse", Song: "Hysteria"},
		{Group: "Coldplay", Song: "Fix You"},
		{ID: 0, Group: "Linkin Park", Song: "Numb"},
		{ID: 1, Group: "Linkin Park", Song: "Numb"},
		{ID: 2, Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{Group: "Linkin Park", Song: "Numb"},
		{ID: 20, Group: "Linkin Park", Song: "Numb"},
		{ID: 21, Group: "Linkin Park", Song: "Numb"},
		{ID: 22, Group: "Linkin Park", Song: "Numb"},
		// Можно добавить больше песен для примера
	}

	// Фильтрация по названию группы или песни (если параметры заданы)
	filteredSongs := []models.Song{}
	for _, s := range allSongs {
		if (group == "" || s.Group == group) && (song == "" || s.Song == song) {
			filteredSongs = append(filteredSongs, s)
		}
	}

	// Общее количество песен после фильтрации
	total := len(filteredSongs)

	// Определение начального и конечного индексов для пагинации
	start := (page - 1) * limit
	end := start + limit

	// Ограничиваем диапазон индексов, чтобы не выйти за пределы массива
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	// Возвращаем только нужные элементы согласно пагинации
	paginatedSongs := filteredSongs[start:end]

	// Возвращаем отфильтрованные и отсортированные песни, общее количество песен и ошибку (если есть)
	return paginatedSongs, total, nil
}
