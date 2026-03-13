package embedder

import (
	"errors"
	"time"
)

var (
	ErrServerNotFound   = errors.New("ollama server not found")
	ErrServerUnexpected = errors.New("ollama server does not response as expected")
	ErrModelNotFound    = errors.New("model not found")
	ErrInvalidURL       = errors.New("invalid url")
	ErrEmptyResponse    = errors.New("empty response from embedder")
)

const (
	defaultTimeout = 10 * time.Second
)
