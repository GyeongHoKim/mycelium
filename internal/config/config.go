// Package config provides types and utilities for parsing and loading
// application configuration.
package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	// FormatModeSection formats output as sections.
	FormatModeSection FormatMode = "section"
	// FormatModeFrontMatter formats output with YAML front matter.
	FormatModeFrontMatter FormatMode = "frontmatter"
)

// Config holds the root configuration for the application.
type Config struct {
	Vault      Vault      `toml:"vault"`
	Embedding  Embedding  `toml:"embedding"`
	Similarity Similarity `toml:"similarity"`
	Output     Output     `toml:"output"`
}

// Vault configures the path to the vault (e.g. Obsidian vault).
type Vault struct {
	Path string `toml:"path"`
}

// Embedding configures the embedding model and Ollama endpoint.
type Embedding struct {
	Model  string `toml:"model"`
	Ollama string `toml:"ollama"`
}

// Similarity configures similarity search (top-k and threshold).
type Similarity struct {
	TopK      int     `toml:"top_k"`
	Threshold float64 `toml:"threshold"`
}

// Output configures output format (section or frontmatter).
type Output struct {
	Format FormatMode `toml:"format"`
}

// FormatMode is the output format: "section" or "frontmatter".
type FormatMode string

// UnmarshalText : encoding.TextUnmarshaler 인터페이스를 구현해주면 toml 라이브러리가 내부에서 호출함.
func (m *FormatMode) UnmarshalText(text []byte) error {
	switch FormatMode(text) {
	case FormatModeSection, FormatModeFrontMatter:
		*m = FormatMode(text)
		return nil
	default:
		return fmt.Errorf("invalid output mode %q: must be section or frontmatter", text)
	}
}

// Load reads configuration from the given path (defaults to "config.toml" if empty).
func Load(path string) (*Config, error) {
	if path == "" {
		path = "config.toml"
	}
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
