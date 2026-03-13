// Package embedder provides Embed, EmbedBatch
package embedder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/gyeonghokim/mycelium/internal/config"
)

type OllamaEmbedder struct {
	client  *http.Client
	baseURL url.URL
	model   string
}

func NewOllamaEmbedder(cfg *config.Config) (*OllamaEmbedder, error) {
	baseURL, err := url.Parse(cfg.Embedding.Ollama)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}
	ollamaEmbedder := &OllamaEmbedder{
		client:  &http.Client{Timeout: defaultTimeout},
		baseURL: *baseURL,
		model:   cfg.Embedding.Model,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err = ollamaEmbedder.Ping(ctx); err != nil {
		return nil, fmt.Errorf("initialization failed: %w", err)
	}

	if err = ollamaEmbedder.CheckModel(ctx); err != nil {
		return nil, fmt.Errorf("initialization failed: %w", err)
	}

	return ollamaEmbedder, nil
}

func (o *OllamaEmbedder) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.baseURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrServerNotFound, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status code %d", ErrServerUnexpected, resp.StatusCode)
	}
	return nil
}

func (o *OllamaEmbedder) CheckModel(ctx context.Context) error {
	u := o.baseURL
	u.Path = "/api/tags"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %w: %w", ErrServerUnexpected, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch tags: %w: status code %d", ErrServerUnexpected, resp.StatusCode)
	}

	var result TagRDO
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode tags response: %w", err)
	}

	for _, m := range result.Models {
		if m.Name == o.model {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrModelNotFound, o.model)
}

func (o *OllamaEmbedder) Embed(_ context.Context, _ string) ([]float32, error) {
	// TODO:
	return nil, errors.New("not implemented")
}

func (o *OllamaEmbedder) EmbedBatch(_ context.Context, _ []string) ([][]float32, error) {
	// TODO:
	return nil, errors.New("not implemented")
}
