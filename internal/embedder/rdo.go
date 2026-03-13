package embedder

type TagRDO struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}
