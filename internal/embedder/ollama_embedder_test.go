package embedder_test

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gyeonghokim/mycelium/internal/config"
	"github.com/gyeonghokim/mycelium/internal/embedder"
)

func setupTestEmbedder(t *testing.T) (*embedder.OllamaEmbedder, *httpmock.MockTransport) {
	t.Helper()
	cfg := &config.Config{
		Embedding: config.Embedding{
			Ollama: "http://localhost:11434/",
			Model:  "qwen3-embedding",
		},
	}
	model := "qwen3-embedding"
	mt := httpmock.NewMockTransport()
	mt.RegisterResponder("GET", "http://localhost:11434/", httpmock.NewStringResponder(200, ""))
	mt.RegisterResponder(
		"GET",
		"http://localhost:11434/api/tags",
		httpmock.NewJsonResponderOrPanic(200, embedder.TagRDO{Models: []embedder.Model{{Size: 1, Name: model}}}),
	)
	client := &http.Client{Transport: mt}

	newEmbedder, err := embedder.NewOllamaEmbedder(cfg, embedder.WithClient(client))
	require.NoError(t, err)

	return newEmbedder, mt
}

func TestOllamaEmbedder_PingSuccess(t *testing.T) {
	t.Parallel()
	newEmbedder, _ := setupTestEmbedder(t)
	err := newEmbedder.Ping(t.Context())
	assert.NoError(t, err)
}

func TestOllamaEmbedder_Embed(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		e, mt := setupTestEmbedder(t)

		expectedEmbedding := embedder.Embedding{0.1, 0.2, 0.3}
		mt.RegisterResponder("POST", "http://localhost:11434/api/embed",
			httpmock.NewJsonResponderOrPanic(200, embedder.EmbedRDO{
				Embeddings: []embedder.Embedding{expectedEmbedding},
			}))

		vec, err := e.Embed(t.Context(), "hello world")
		require.NoError(t, err)
		assert.Equal(t, expectedEmbedding, vec)
	})

	t.Run("server error", func(t *testing.T) {
		t.Parallel()
		e, mt := setupTestEmbedder(t)

		mt.RegisterResponder("POST", "http://localhost:11434/api/embed",
			httpmock.NewStringResponder(500, "internal server error"))

		vec, err := e.Embed(t.Context(), "hello world")
		require.Error(t, err)
		assert.Nil(t, vec)
	})

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()
		e, mt := setupTestEmbedder(t)

		// Ollama returning empty embeddings list
		mt.RegisterResponder("POST", "http://localhost:11434/api/embed",
			httpmock.NewJsonResponderOrPanic(200, embedder.EmbedRDO{
				Embeddings: []embedder.Embedding{},
			}))

		vec, err := e.Embed(t.Context(), "hello world")
		require.ErrorIs(t, err, embedder.ErrEmptyResponse)
		assert.Nil(t, vec)
	})
}

func TestOllamaEmbedder_EmbedBatch(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		e, mt := setupTestEmbedder(t)

		expectedEmbeddings := []embedder.Embedding{
			{0.1, 0.2},
			{0.3, 0.4},
		}
		mt.RegisterResponder("POST", "http://localhost:11434/api/embed",
			httpmock.NewJsonResponderOrPanic(200, embedder.EmbedRDO{
				Embeddings: expectedEmbeddings,
			}))

		vecs, err := e.EmbedBatch(t.Context(), []string{"hello", "world"})
		require.NoError(t, err)
		assert.Len(t, vecs, 2)
		assert.Equal(t, expectedEmbeddings[0], vecs[0])
		assert.Equal(t, expectedEmbeddings[1], vecs[1])
	})

	t.Run("server error", func(t *testing.T) {
		t.Parallel()
		e, mt := setupTestEmbedder(t)

		mt.RegisterResponder("POST", "http://localhost:11434/api/embed",
			httpmock.NewStringResponder(500, "internal server error"))

		vecs, err := e.EmbedBatch(t.Context(), []string{"hello", "world"})
		require.Error(t, err)
		assert.Nil(t, vecs)
	})
}
