package ai

import (
	"context"
	"errors"
	"os"

	"github.com/pgvector/pgvector-go"
	"google.golang.org/genai"
)

const (
	embeddingModel = "gemini-embedding-001"
	embeddingDim   = 3072
)

// NewClient creates a reusable genai client.
func NewClient(ctx context.Context) (*genai.Client, error) {
	if os.Getenv("GOOGLE_API_KEY") == "" {
		return nil, errors.New("embedding: GOOGLE_API_KEY not set")
	}
	c, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, errors.New("new client: " + err.Error())
	}
	return c, nil
}

// Embedding generates an embedding vector for text using a provided client.
func Embedding(ctx context.Context, client *genai.Client, text []string) (pgvector.Vector, error) {
	if len(text) == 0 {
		return pgvector.Vector{}, errors.New("embedding: empty text")
	}

	contents := []*genai.Content{}
	for _, t := range text {
		contents = append(contents, genai.NewContentFromText(t, genai.RoleUser))
	}

	res, err := client.Models.EmbedContent(
		ctx,
		embeddingModel,
		contents,
		&genai.EmbedContentConfig{OutputDimensionality: func() *int32 { v := int32(embeddingDim); return &v }()},
	)
	if err != nil {
		return pgvector.Vector{}, errors.New("embedding: " + err.Error())
	}
	if len(res.Embeddings) == 0 || len(res.Embeddings[0].Values) == 0 {
		return pgvector.Vector{}, errors.New("embedding: empty result")
	}

	raw := res.Embeddings[0].Values
	if len(raw) != embeddingDim {
		return pgvector.Vector{}, errors.New("embedding: invalid dimension")
	}

	out := make([]float32, len(raw))
	for i, v := range raw {
		out[i] = float32(v)
	}
	return pgvector.NewVector(out), nil
}
