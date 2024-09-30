package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
	"github.com/pkg/errors"
)

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) GetSongByID(id int) (models.Song, error) {
	var song models.Song
	query := `
		SELECT 
			id, 
			group_name, 
			title, 
			lyrics, 
			release_date, 
			link 
		FROM songs 
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(&song.ID, &song.Group, &song.Title, &song.Lyrics, &song.ReleaseDate, &song.Link)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Song{}, fmt.Errorf("no song found with ID %d", id)
		}
		return models.Song{}, err
	}
	return song, nil
}

func (r *SongRepository) GetSongsByAnyField(group, title string) ([]models.Song, error) {
	query := `SELECT id, group_name, title, release_date, lyrics, link FROM songs WHERE 1=1`
	args := []interface{}{}

	if group != "" {
		query += " AND group_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, group)
	}
	if title != "" {
		query += " AND title = $" + strconv.Itoa(len(args)+1)
		args = append(args, title)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query execution error")
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		var releaseDate, lyrics, link sql.NullString

		if err := rows.Scan(&song.ID, &song.Group, &song.Title, &releaseDate, &lyrics, &link); err != nil {
			return nil, errors.Wrap(err, "row scanning error")
		}

		if releaseDate.Valid {
			song.ReleaseDate = &releaseDate.String
		} else {
			song.ReleaseDate = nil
		}

		if lyrics.Valid {
			song.Lyrics = &lyrics.String
		} else {
			song.Lyrics = nil
		}

		if link.Valid {
			song.Link = &link.String
		} else {
			song.Link = nil
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (r *SongRepository) DeleteSongByID(id int) error {
	query := "DELETE FROM songs WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "rows affected")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no song found with ID %d", id)
	}

	return nil
}

func (r *SongRepository) UpdateSongByID(id int, updateRequest *models.UpdateSongRequest) error {
	query := "UPDATE songs SET "
	var params []interface{}
	var setClauses []string
	paramCount := 1
	if updateRequest.Group != nil {
		setClauses = append(setClauses, fmt.Sprintf("group_name = $%d", paramCount))
		params = append(params, *updateRequest.Group)
		paramCount++
	}
	if updateRequest.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", paramCount))
		params = append(params, *updateRequest.Title)
		paramCount++
	}
	if updateRequest.Lyrics != nil {
		setClauses = append(setClauses, fmt.Sprintf("lyrics = $%d", paramCount))
		params = append(params, *updateRequest.Lyrics)
		paramCount++
	}
	if updateRequest.ReleaseDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("release_date = $%d", paramCount))
		params = append(params, *updateRequest.ReleaseDate)
		paramCount++
	}
	if updateRequest.Link != nil {
		setClauses = append(setClauses, fmt.Sprintf("link = $%d", paramCount))
		params = append(params, *updateRequest.Link)
		paramCount++
	}

	if len(setClauses) == 0 {
		return errors.New("no fields provided for update")
	}

	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	result, err := r.db.Exec(query, params...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "rows affected")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no song found with ID %d", id)
	}

	return nil
}

func (r *SongRepository) AddSong(group, song string, songDetail *models.SongDetail) (int, error) {
	formattedDate, err := convertDate(songDetail.ReleaseDate)
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO songs (group_name, title, lyrics, release_date, link) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id int
	err = r.db.QueryRow(query, group, song, songDetail.Text, formattedDate, songDetail.Link).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func convertDate(dateStr string) (string, error) {
	parsedDate, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %w", err)
	}

	return parsedDate.Format("2006-01-02"), nil
}
