package models

// Song представляет данные для добавления новой песни
type Song struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// SongDetail представляет информацию, получаемую от внешнего API
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
