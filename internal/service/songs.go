package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
	"github.com/KarmaBeLike/SongLibrary/internal/repository"
)

type SongService struct {
	songRepo       *repository.SongRepository
	externalAPIURL string
}

func NewSongService(songRepo *repository.SongRepository, externalAPIURL string) *SongService {
	return &SongService{
		songRepo:       songRepo,
		externalAPIURL: externalAPIURL,
	}
}

func (s *SongService) GetSongs(group, song string, page, limit int) ([]models.Song, models.Pagination, error) {
	songs, err := s.songRepo.GetSongsByAnyField(group, song)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	total := len(songs)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedSongs := songs[start:end]

	pagination := models.Pagination{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return paginatedSongs, pagination, nil
}

func (s *SongService) GetPaginatedSongLyrics(id int, page, limit int) (*models.SongVerses, error) {
	song, err := s.songRepo.GetSongByID(id)
	if err != nil {
		return nil, err
	}

	if song.Lyrics == nil {
		return nil, fmt.Errorf("lyrics not found for song ID %d", id)
	}

	verses := strings.Split(*song.Lyrics, "\n\n")

	for i := range verses {
		verses[i] = strings.TrimSpace(verses[i])
	}

	totalVerses := len(verses)

	start := (page - 1) * limit
	end := start + limit

	if start >= totalVerses || start < 0 {
		return nil, fmt.Errorf("page out of range")
	}

	if end > totalVerses {
		end = totalVerses
	}

	paginatedVerses := verses[start:end]

	response := &models.SongVerses{
		ID:     song.ID,
		Group:  song.Group,
		Title:  song.Title,
		Verses: paginatedVerses,
	}

	return response, nil
}

func (s *SongService) DeleteSongByID(id int) error {
	err := s.songRepo.DeleteSongByID(id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}
	return nil
}

func (s *SongService) UpdateSongByID(id int, updateRequest *models.UpdateSongRequest) error {
	return s.songRepo.UpdateSongByID(id, updateRequest)
}

func (s *SongService) AddSong(newSong models.NewSongRequest) (int, error) {
	slog.Info("Adding new song", slog.String("group", newSong.Group), slog.String("song", newSong.Song))

	requestURL := fmt.Sprintf("%s?group=%s&song=%s", s.externalAPIURL, url.QueryEscape(newSong.Group), url.QueryEscape(newSong.Song))

	slog.Debug("Calling external API", slog.String("url", requestURL))

	resp, err := http.Get(requestURL)
	if err != nil {
		slog.Error("Failed to call external API", slog.Any("error", err))
		return 0, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()

	slog.Debug("External API response", slog.String("status", resp.Status))

	if resp.StatusCode != http.StatusOK {
		slog.Warn("External API returned non-OK status", slog.String("status", resp.Status))
		return 0, fmt.Errorf("external API returned status: %s", resp.Status)
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		slog.Error("Failed to decode external API response", slog.Any("error", err))
		return 0, fmt.Errorf("failed to decode external API response: %w", err)
	}

	slog.Info("Successfully retrieved song details from external API", slog.Any("songDetail", songDetail))

	id, err := s.songRepo.AddSong(newSong.Group, newSong.Song, &songDetail)
	if err != nil {
		slog.Error("Failed to add song to the database", slog.Any("error", err))
		return 0, err
	}
	slog.Info("Successfully added song to the database", slog.Int("song_id", id))

	return id, nil
}
