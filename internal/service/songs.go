package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/KarmaBeLike/SongLibrary/internal/api"
	"github.com/KarmaBeLike/SongLibrary/internal/models"
	"github.com/KarmaBeLike/SongLibrary/internal/repository"
	"github.com/KarmaBeLike/SongLibrary/pkg/validation"
)

type SongService struct {
	songRepo    *repository.SongRepository
	externalAPI *api.ExternalAPI
}

func NewSongService(songRepo *repository.SongRepository, externalAPI *api.ExternalAPI) *SongService {
	return &SongService{
		songRepo:    songRepo,
		externalAPI: externalAPI,
	}
}

func (s *SongService) GetSongs(ctx context.Context, filter models.SongFilter) ([]models.Song, models.Pagination, error) {
	slog.Debug("Fetching songs", slog.Any("filter", filter))

	// Получение песен через репозиторий с фильтром
	songs, err := s.songRepo.GetSongsByFilter(ctx, filter)
	if err != nil {
		slog.Error("Error fetching songs", slog.Any("filter", filter), slog.Any("error", err))
		return nil, models.Pagination{}, err
	}

	// Пагинация
	total := len(songs)
	start := (filter.Page - 1) * filter.Limit
	end := start + filter.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedSongs := songs[start:end]

	pagination := models.Pagination{
		Limit: filter.Limit,
		Page:  filter.Page,
		Total: total,
	}

	slog.Info("Songs fetched successfully", slog.Int("total", total), slog.Int("returned_count", len(paginatedSongs)))

	return paginatedSongs, pagination, nil
}

func (s *SongService) GetPaginatedSongLyrics(ctx context.Context, id int, page, limit int) (*models.SongVerses, error) {
	slog.Debug("Fetching song lyrics", slog.Int("song_id", id), slog.Int("page", page), slog.Int("limit", limit))

	song, err := s.songRepo.GetSongByID(ctx, id)
	if err != nil {
		slog.Error("Error fetching song by ID", slog.Int("song_id", id), slog.Any("error", err))
		return nil, err
	}

	if song.Text == nil {
		slog.Warn("Text not found for song", slog.Int("song_id", id))
		return nil, fmt.Errorf("lyrics not found for song ID %d", id)
	}

	verses := strings.Split(*song.Text, "\n\n")

	for i := range verses {
		verses[i] = strings.TrimSpace(verses[i])
	}

	totalVerses := len(verses)

	start := (page - 1) * limit
	end := start + limit

	if start >= totalVerses || start < 0 {
		slog.Warn("Page out of range", slog.Int("page", page), slog.Int("total_verses", totalVerses))
		return nil, fmt.Errorf("page out of range")
	}

	if end > totalVerses {
		end = totalVerses
	}

	paginatedVerses := verses[start:end]

	response := &models.SongVerses{
		ID:     song.ID,
		Group:  song.Group,
		Song:   song.Song,
		Verses: paginatedVerses,
		Total:  len(verses),
	}

	slog.Info("Text fetched successfully", slog.Int("song_id", id), slog.Int("returned_verses", len(paginatedVerses)))

	return response, nil
}

func (s *SongService) DeleteSongByID(ctx context.Context, id int) error {
	slog.Debug("Deleting song by ID", slog.Int("song_id", id))

	err := s.songRepo.DeleteSongByID(ctx, id)
	if err != nil {
		slog.Error("Error deleting song", slog.Int("song_id", id), slog.Any("error", err))
		return fmt.Errorf("failed to delete song: %w", err)
	}
	slog.Info("Song deleted successfully", slog.Int("song_id", id))
	return nil
}

func (s *SongService) UpdateSongByID(ctx context.Context, id int, updateRequest *models.UpdateSongRequest) error {
	slog.Debug("Updating song by ID", slog.Int("song_id", id), slog.Any("updateRequest", updateRequest))
	err := s.songRepo.UpdateSongByID(ctx, id, updateRequest)
	if err != nil {
		slog.Error("Error updating song", slog.Int("song_id", id), slog.Any("error", err))
		return err
	}

	slog.Info("Song updated successfully", slog.Int("song_id", id))
	return nil
}

func (s *SongService) AddSong(ctx context.Context, newSong models.NewSongRequest) (int, error) {
	slog.Info("Adding new song", slog.String("group", newSong.Group), slog.String("song", newSong.Song))

	// 2. Получить детали песни из внешнего API
	songDetail, err := s.externalAPI.FetchSongDetail(ctx, newSong.Group, newSong.Song)
	if err != nil {
		slog.Error("Failed to fetch song detail", slog.Any("error", err))
		return 0, fmt.Errorf("failed to fetch song detail: %w", err)
	}

	// 3. Проверить текст песни, если нужно
	if err := validation.ValidateSongText(songDetail.Text); err != nil {
		slog.Error("Song text validation failed", slog.Any("error", err))
		return 0, fmt.Errorf("song text validation error: %w", err)
	}

	// 4. Добавить песню в базу данных
	songID, err := s.songRepo.AddSong(ctx, newSong.Group, newSong.Song, songDetail)
	if err != nil {
		slog.Error("Failed to add song to the database", slog.Any("error", err))
		return 0, err
	}

	slog.Info("Successfully added song to the database", slog.Int("songID", songID))
	return songID, nil
}
