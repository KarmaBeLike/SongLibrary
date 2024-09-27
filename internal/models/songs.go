package models

// Song представляет данные для добавления новой песни
type Song struct {
	ID     int      `json:"id"`
	Group  string   `json:"group"`
	Song   string   `json:"song"`
	Limit  int      `json:"limit"`
	Page   int      `json:"page"`
	Total  int      `json:"total"`
	Verses []string `json:"verses"`
}

// SongDetail представляет информацию, получаемую от внешнего API
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
