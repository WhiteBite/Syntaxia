package embeddings

import (
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
		delta    float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
			delta:    0.001,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{-1, 0, 0},
			expected: -1.0,
			delta:    0.001,
		},
		{
			name:     "similar vectors",
			a:        []float32{1, 1, 0},
			b:        []float32{1, 0, 0},
			expected: 0.707,
			delta:    0.01,
		},
		{
			name:     "empty vectors",
			a:        []float32{},
			b:        []float32{},
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "different lengths",
			a:        []float32{1, 0},
			b:        []float32{1, 0, 0},
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "zero vector",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 0.0,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			diff := result - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.delta {
				t.Errorf("cosineSimilarity() = %v, want %v (±%v)", result, tt.expected, tt.delta)
			}
		})
	}
}

func TestEncodeDecodeEmbedding(t *testing.T) {
	original := []float32{0.1, 0.2, 0.3, 0.4, 0.5}

	encoded, err := encodeEmbedding(original)
	if err != nil {
		t.Fatalf("encodeEmbedding failed: %v", err)
	}

	decoded, err := decodeEmbedding(encoded)
	if err != nil {
		t.Fatalf("decodeEmbedding failed: %v", err)
	}

	if len(decoded) != len(original) {
		t.Errorf("length mismatch: got %d, want %d", len(decoded), len(original))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("value mismatch at %d: got %v, want %v", i, decoded[i], original[i])
		}
	}
}

func TestEncodeDecodeEmbedding_Empty(t *testing.T) {
	original := []float32{}

	encoded, err := encodeEmbedding(original)
	if err != nil {
		t.Fatalf("encodeEmbedding failed: %v", err)
	}

	decoded, err := decodeEmbedding(encoded)
	if err != nil {
		t.Fatalf("decodeEmbedding failed: %v", err)
	}

	if len(decoded) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(decoded))
	}
}

func TestEncodeDecodeEmbedding_Large(t *testing.T) {
	// Test with typical embedding size (1536 for OpenAI ada-002)
	original := make([]float32, 1536)
	for i := range original {
		original[i] = float32(i) / 1536.0
	}

	encoded, err := encodeEmbedding(original)
	if err != nil {
		t.Fatalf("encodeEmbedding failed: %v", err)
	}

	decoded, err := decodeEmbedding(encoded)
	if err != nil {
		t.Fatalf("decodeEmbedding failed: %v", err)
	}

	if len(decoded) != len(original) {
		t.Errorf("length mismatch: got %d, want %d", len(decoded), len(original))
	}
}

func TestDecodeEmbedding_InvalidJSON(t *testing.T) {
	_, err := decodeEmbedding([]byte("invalid json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCosineSimilarity_Normalized(t *testing.T) {
	// Test with normalized vectors (common in embeddings)
	a := []float32{0.6, 0.8, 0} // normalized: sqrt(0.36 + 0.64) = 1
	b := []float32{0.8, 0.6, 0} // normalized: sqrt(0.64 + 0.36) = 1

	result := cosineSimilarity(a, b)
	expected := float32(0.96) // 0.6*0.8 + 0.8*0.6 = 0.96

	diff := result - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.01 {
		t.Errorf("cosineSimilarity() = %v, want %v", result, expected)
	}
}
