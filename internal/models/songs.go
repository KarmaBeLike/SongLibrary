package models

type Song struct {
	ID          int     `json:"song_id"`
	Group       string  `json:"group"`
	Title       string  `json:"song"`
	Lyrics      *string `json:"text"`
	ReleaseDate *string `json:"releaseDate"`
	Link        *string `json:"link"`
}
type Group struct {
	ID   int    `json:"group_id"`
	Name string `json:"name"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Total int `json:"total"`
}

type SongFilter struct {
	Group       string `json:"group"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	ReleaseDate string `json:"releaseDate"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type SongVerses struct {
	ID     int      `json:"id"`
	Group  string   `json:"group"`
	Title  string   `json:"song"`
	Verses []string `json:"verses"`
}

type UpdateSongRequest struct {
	Group       *string `json:"group,omitempty"`
	Title       *string `json:"title,omitempty"`
	Lyrics      *string `json:"lyrics,omitempty"`
	ReleaseDate *string `json:"releaseDate,omitempty"`
	Link        *string `json:"link,omitempty"`
}

type NewSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
