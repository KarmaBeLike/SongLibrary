package repository

import (
	"context"
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

func (r *SongRepository) GetSongByID(ctx context.Context, id int) (models.Song, error) {
	slog.Debug("Fetching song by ID", slog.Int("id", id))

	var song models.Song
	query := `
		SELECT 
			song_id, 
			group_name, 
			song, 
			lyrics, 
			release_date, 
			link 
		FROM songs 
		WHERE song_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&song.ID, &song.Group, &song.Song, &song.Text, &song.ReleaseDate, &song.Link)
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

func (r *SongRepository) GetSongsByFilter(ctx context.Context, filter models.SongFilter) ([]models.Song, error) {
	query := `SELECT  song_id,group_name, song, lyrics, release_date, link FROM songs WHERE 1=1`
	args := []interface{}{}
	paramIndex := 1

	if filter.Group != "" {
		query += fmt.Sprintf(" AND group_name = $%d", paramIndex)
		args = append(args, filter.Group)
		paramIndex++
	}

	if filter.Song != "" {
		query += fmt.Sprintf(" AND song ILIKE $%d", paramIndex)
		args = append(args, "%"+filter.Song+"%")
		paramIndex++
	}

	if filter.Text != "" {
		query += fmt.Sprintf(" AND lyrics ILIKE $%d", paramIndex)
		args = append(args, "%"+filter.Text+"%")
		paramIndex++
	}

	if filter.Link != "" {
		query += fmt.Sprintf(" AND link = $%d", paramIndex)
		args = append(args, filter.Link)
		paramIndex++
	}

	if filter.ReleaseDate != "" {
		query += fmt.Sprintf(" AND release_date = $%d", paramIndex)
		args = append(args, filter.ReleaseDate)
		paramIndex++
	}

	slog.Debug("Executing query", slog.String("query", query), slog.Any("args", args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.Text, &song.ReleaseDate, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (r *SongRepository) DeleteSongByID(ctx context.Context, id int) error {
	slog.Debug("Deleting song by ID", slog.Int("id", id))

	query := "DELETE FROM songs WHERE song_id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *SongRepository) UpdateSongByID(ctx context.Context, id int, updateRequest *models.UpdateSongRequest) error {
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
	if updateRequest.Song != nil {
		setClauses = append(setClauses, fmt.Sprintf("song = $%d", paramCount))
		params = append(params, *updateRequest.Song)
		paramCount++
	}
	if updateRequest.Text != nil {
		setClauses = append(setClauses, fmt.Sprintf("lyrics = $%d", paramCount))
		params = append(params, *updateRequest.Text)
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

	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE song_id = $%d", paramCount)
	params = append(params, id)

	result, err := r.db.ExecContext(ctx, query, params...)
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

func (r *SongRepository) GetOrCreateGroup(ctx context.Context, groupName string) (int, error) {
	var groupID int

	// Проверить, существует ли группа
	queryCheck := "SELECT group_id FROM groups WHERE name = $1"
	err := r.db.QueryRowContext(ctx, queryCheck, groupName).Scan(&groupID)
	if err == nil {
		return groupID, nil // Группа найдена, вернуть её ID
	}

	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to check group: %w", err)
	}

	// Если группы нет, добавить её
	queryInsert := "INSERT INTO groups (name) VALUES ($1) RETURNING group_id"
	err = r.db.QueryRowContext(ctx, queryInsert, groupName).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert group: %w", err)
	}

	return groupID, nil
}

func (r *SongRepository) AddSong(ctx context.Context, group, song string, songDetail *models.SongDetail) (int, error) {
	slog.Debug("Adding new song", slog.String("group", group), slog.String("song", song))

	formattedDate, err := convertDate(songDetail.ReleaseDate)
	if err != nil {
		slog.Error("Error converting date", slog.Any("error", err))
		return 0, err
	}

	groupID, err := r.GetOrCreateGroup(ctx, group)
	if err != nil {
		slog.Error("Error getting or creating group", slog.Any("error", err))
		return 0, err
	}

	slog.Debug("Group ID obtained", slog.Int("groupID", groupID))

	// Вставить песню
	query := "INSERT INTO songs (group_name, song, lyrics, release_date, link,group_id) VALUES ($1, $2, $3, $4, $5,$6) RETURNING song_id"
	var songID int
	err = r.db.QueryRowContext(ctx, query, group, song, songDetail.Text, formattedDate, songDetail.Link, groupID).Scan(&songID)
	if err != nil {
		slog.Error("Error inserting new song", slog.Any("error", err))
		return 0, err
	}

	slog.Info("Song added successfully", slog.Int("songID", songID))
	return songID, nil
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
