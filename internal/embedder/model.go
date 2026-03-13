package embedder

import "math"

type Embedding []float32

// L2Norm 은 테스트 코드 혹은 모킹 Embedder 만들 때 사용.
func L2Norm(v []float32) float64 {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	return math.Sqrt(sum)
}

// CosineSimilarity 은 테스트 코드 혹은 모킹 Embedder 만들 때 사용.
func CosineSimilarity(a, b []float32) float32 {
	var dot float32
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot
}
