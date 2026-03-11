package config_test

import (
	"testing"

	"github.com/gyeonghokim/mycelium/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestFormatModeUnmarshalText(t *testing.T) {
	t.Parallel()

	formatModeTests := []struct {
		name    string
		input   []byte
		want    config.FormatMode
		wantErr bool
	}{
		{"section", []byte("section"), config.FormatModeSection, false},
		{"frontmatter", []byte("frontmatter"), config.FormatModeFrontMatter, false},
		{"invalid", []byte("invalid"), "", true},
	}

	for _, tt := range formatModeTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var m config.FormatMode
			err := m.UnmarshalText(tt.input)
			assert.Equal(t, tt.wantErr, err != nil, "error expectation")
			if !tt.wantErr {
				assert.Equal(t, tt.want, m, "unmarshaled value")
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Parallel()

	loadTests := []struct {
		name    string
		path    string
		want    *config.Config
		wantErr bool
	}{
		{
			name: "valid toml",
			path: "testdata/valid.toml",
			want: &config.Config{
				Vault: config.Vault{Path: "/path/to/vault"},
				Embedding: config.Embedding{
					Model:  "nomic-embed-text",
					Ollama: "http://localhost:11434",
				},
				Similarity: config.Similarity{TopK: 5, Threshold: 0.7},
				Output:     config.Output{Format: config.FormatModeSection},
			},
			wantErr: false,
		},
		{
			name:    "invalid syntax",
			path:    "testdata/invalid_syntax.toml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format",
			path:    "testdata/invalid_format.toml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "file not found",
			path:    "testdata/nonexistent.toml",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range loadTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := config.Load(tt.path)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.want, cfg)
			}
		})
	}
}
