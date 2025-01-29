package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/KarmaBeLike/SongLibrary/internal/models"
)

type ExternalAPI struct {
	baseURL string
	client  *http.Client
}

func NewExternalAPI(baseURL string) *ExternalAPI {
	return &ExternalAPI{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ExternalAPI) FetchSongDetail(ctx context.Context, group, song string) (*models.SongDetail, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	q := u.Query()
	q.Set("group", group)
	q.Set("song", song)
	u.RawQuery = q.Encode()

	slog.Debug("Fetching song detail", slog.String("url", u.String()))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch song detail: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &songDetail, nil
}
