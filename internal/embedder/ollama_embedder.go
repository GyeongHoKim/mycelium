// Package embedder provides Embed, EmbedBatch
package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gyeonghokim/mycelium/internal/config"
)

type OllamaEmbedder struct {
	client  *http.Client
	baseURL url.URL
	model   string
}

type Option func(*OllamaEmbedder)

func WithClient(client *http.Client) Option {
	return func(o *OllamaEmbedder) {
		o.client = client
	}
}

func NewOllamaEmbedder(cfg *config.Config, opts ...Option) (*OllamaEmbedder, error) {
	baseURL, err := url.Parse(cfg.Embedding.Ollama)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}
	ollamaEmbedder := &OllamaEmbedder{
		client:  &http.Client{Timeout: defaultTimeout},
		baseURL: *baseURL,
		model:   cfg.Embedding.Model,
	}

	for _, opt := range opts {
		opt(ollamaEmbedder)
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

func (o *OllamaEmbedder) Embed(ctx context.Context, text string) (Embedding, error) {
	dto := &EmbedRequestDTO{
		Model: o.model,
		Input: text,
	}
	u := o.baseURL
	u.Path = "/api/embed"

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(dto); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmbedRDO
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embeddings[0], nil
}

func (o *OllamaEmbedder) EmbedBatch(ctx context.Context, texts []string) ([]Embedding, error) {
	dto := &EmbedBatchRequestDTO{
		Model: o.model,
		Input: texts,
	}
	u := o.baseURL
	u.Path = "/api/embed"

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(dto); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmbedRDO
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embeddings, nil
}
