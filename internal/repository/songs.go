package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
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
	slog.Debug("Fetching song by ID", slog.Int("id", id))

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
			slog.Info("No song found with ID", slog.Int("id", id))
			return models.Song{}, fmt.Errorf("no song found with ID %d", id)
		}
		slog.Error("Error fetching song by ID", slog.Any("error", err))
		return models.Song{}, err
	}
	slog.Info("Song fetched successfully", slog.Int("id", id))
	return song, nil
}

func (r *SongRepository) GetSongsByFilter(filter models.SongFilter) ([]models.Song, error) {
	query := `SELECT  song_id,group_name, title, lyrics, release_date, link FROM songs WHERE 1=1`
	args := []interface{}{}
	paramIndex := 1

	if filter.Group != "" {
		query += fmt.Sprintf(" AND group_name = $%d", paramIndex)
		args = append(args, filter.Group)
		paramIndex++
	}
	if filter.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", paramIndex)
		args = append(args, "%"+filter.Title+"%")
		paramIndex++
	}
	if filter.Text != "" {
		query += fmt.Sprintf(" AND lyrics ILIKE $%d", paramIndex)
		args = append(args, "%"+filter.Text+"%")
		paramIndex++
	}
	if filter.ReleaseDate != "" {
		query += fmt.Sprintf(" AND release_date = $%d", paramIndex)
		args = append(args, filter.ReleaseDate)
		paramIndex++
	}

	slog.Debug("Executing query", slog.String("query", query), slog.Any("args", args))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Title, &song.Lyrics, &song.ReleaseDate, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (r *SongRepository) DeleteSongByID(id int) error {
	slog.Debug("Deleting song by ID", slog.Int("id", id))

	query := "DELETE FROM songs WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		slog.Error("Error executing delete query", slog.Any("error", err))
		return errors.Wrap(err, "execute query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Error getting rows affected", slog.Any("error", err))
		return errors.Wrap(err, "rows affected")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no song found with ID %d", id)
	}
	slog.Info("Song deleted successfully", slog.Int("id", id))
	return nil
}

func (r *SongRepository) UpdateSongByID(id int, updateRequest *models.UpdateSongRequest) error {
	slog.Debug("Updating song by ID", slog.Int("id", id), slog.Any("updateRequest", updateRequest))

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
		slog.Warn("No fields provided for update", slog.Int("id", id))
		return errors.New("no fields provided for update")
	}

	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	result, err := r.db.Exec(query, params...)
	if err != nil {
		slog.Error("Error executing update query", slog.Any("error", err))
		return errors.Wrap(err, "execute query")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Error getting rows affected", slog.Any("error", err))
		return errors.Wrap(err, "rows affected")
	}
	if rowsAffected == 0 {
		slog.Info("No song found to update", slog.Int("id", id))
		return fmt.Errorf("no song found with ID %d", id)
	}

	slog.Info("Song updated successfully", slog.Int("id", id))
	return nil
}

func (r *SongRepository) AddSong(group, song string, songDetail *models.SongDetail) (int, error) {
	slog.Debug("Adding new song", slog.String("group", group), slog.String("title", song))

	formattedDate, err := convertDate(songDetail.ReleaseDate)
	if err != nil {
		slog.Error("Error converting date", slog.Any("error", err))
		return 0, err
	}

	query := "INSERT INTO songs (group_name, title, lyrics, release_date, link) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id int
	err = r.db.QueryRow(query, group, song, songDetail.Text, formattedDate, songDetail.Link).Scan(&id)
	if err != nil {
		slog.Error("Error inserting new song", slog.Any("error", err))
		return 0, err
	}
	slog.Info("Song added successfully", slog.Int("id", id))
	return id, nil
}

func convertDate(dateStr string) (string, error) {
	slog.Debug("Converting date", slog.String("dateStr", dateStr))

	parsedDate, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		slog.Error("Error parsing date", slog.Any("error", err))
		return "", fmt.Errorf("failed to parse date: %w", err)
	}

	return parsedDate.Format("2006-01-02"), nil
}
