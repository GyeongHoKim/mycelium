package embedder

type EmbedRequestDTO struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedBatchRequestDTO struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}
