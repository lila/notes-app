package concept

import (
	"log"
	"note-app/ai"
	"note-app/cache"
	"note-app/note"
	"sort"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type Concept struct {
	Name  string
	Notes []note.Note
}

// generateCacheKey creates a stable, unique key from a slice of notes.
func generateCacheKey(notes []note.Note) string {
	ids := make([]int, len(notes))
	for i, n := range notes {
		ids[i] = n.ID
	}
	sort.Ints(ids) // Sort to ensure stability

	var sb strings.Builder
	for i, id := range ids {
		sb.WriteString(strconv.Itoa(id))
		if i < len(ids)-1 {
			sb.WriteString("-")
		}
	}
	return sb.String()
}

func Cluster(notes []note.Note, force bool) []Concept {
	if len(notes) == 0 {
		return nil
	}

	var concepts []Concept
	used := make(map[int]bool)

	for i := 0; i < len(notes); i++ {
		if used[i] || len(notes[i].Embedding) == 0 {
			continue
		}

		// Start a new cluster
		clusterNotes := []note.Note{notes[i]}
		used[i] = true

		// Find similar notes
		for j := i + 1; j < len(notes); j++ {
			if used[j] || len(notes[j].Embedding) == 0 {
				continue
			}
			similarity := cosineSimilarity(notes[i].Embedding, notes[j].Embedding)
			if similarity > 0.8 {
				clusterNotes = append(clusterNotes, notes[j])
				used[j] = true
			}
		}

		// Check cache first unless a regeneration has been forced.
		cacheKey := generateCacheKey(clusterNotes)
		if !force {
			if name, found := cache.Get(cacheKey); found {
				concepts = append(concepts, Concept{Name: name, Notes: clusterNotes})
				continue
			}
		}

		// If not in cache, generate a name for the cluster
		var contents []string
		for _, n := range clusterNotes {
			contents = append(contents, n.Content)
		}
		name, err := ai.GenerateConceptName(contents)
		if err != nil {
			log.Printf("Failed to generate concept name: %v", err)
			name = "Unnamed Concept"
		}

		// Save to cache
		cache.Set(cacheKey, name)
		concepts = append(concepts, Concept{Name: name, Notes: clusterNotes})
	}

	sort.Slice(concepts, func(i, j int) bool {
		return len(concepts[i].Notes) > len(concepts[j].Notes)
	})

	return concepts
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	vecA := mat.NewVecDense(len(a), float64Array(a))
	vecB := mat.NewVecDense(len(b), float64Array(b))
	dot := mat.Dot(vecA, vecB)
	normA := mat.Norm(vecA, 2)
	normB := mat.Norm(vecB, 2)
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (normA * normB)
}

func float64Array(a []float32) []float64 {
	res := make([]float64, len(a))
	for i, v := range a {
		res[i] = float64(v)
	}
	return res
}
