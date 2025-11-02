package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	embeddingModel *genai.EmbeddingModel
	genModel       *genai.GenerativeModel
)

func init() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("GEMINI_API_KEY not set, AI features will be disabled.")
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create new AI client: %v", err)
	}

	embeddingModel = client.EmbeddingModel("embedding-001")

	iter := client.ListModels(ctx)
	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list models: %v", err)
		}
		for _, method := range m.SupportedGenerationMethods {
			if method == "generateContent" {
				genModel = client.GenerativeModel(m.Name)
				return
			}
		}
	}
	log.Println("No suitable generative model found.")
}

func GenerateEmbedding(text string) ([]float32, error) {
	if embeddingModel == nil {
		return nil, fmt.Errorf("embedding model not initialized")
	}

	ctx := context.Background()
	res, err := embeddingModel.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}
	return res.Embedding.Values, nil
}

func GenerateConceptName(noteContents []string) (string, error) {
	if genModel == nil {
		return "Unnamed Concept", fmt.Errorf("generative model not initialized")
	}

	prompt := fmt.Sprintf(
		"Based on the following notes, what is a short, 2-3 word topic or concept name for them? Notes:\n\n- %s",
		strings.Join(noteContents, "\n- "),
	)

	ctx := context.Background()
	resp, err := genModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "Unnamed Concept", err
	}

	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				return string(txt), nil
			}
		}
	}

	return "Unnamed Concept", nil
}
