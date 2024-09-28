package repository

import "database/sql"

func InsertSong(db *sql.DB, group, song, releaseDate, lyrics, link string) error {
	query := `INSERT INTO songs (group_name, song_title, release_date, lyrics, link) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(query, group, song, releaseDate, lyrics, link)
	if err != nil {
		return err
	}

	return nil
}
