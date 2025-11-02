package note

import (
	"encoding/json"
	"log"
	"time"

	"note-app/ai"
	"note-app/store"
)

type Note struct {
	ID        int
	Content   string
	Created   time.Time
	Embedding []float32
}

var notes []Note
var nextID = 1

func init() {
	data, err := store.Load()
	if err != nil {
		log.Fatalf("Failed to load notes: %v", err)
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &notes); err != nil {
			log.Fatalf("Failed to unmarshal notes: %v", err)
		}
		if len(notes) > 0 {
			nextID = notes[len(notes)-1].ID + 1
		}
	}
}

func Add(content string) {
	embedding, err := ai.GenerateEmbedding(content)
	if err != nil {
		log.Printf("Failed to generate embedding: %v", err)
	}

	notes = append(notes, Note{
		ID:        nextID,
		Content:   content,
		Created:   time.Now(),
		Embedding: embedding,
	})
	nextID++
	if err := store.Save(notes); err != nil {
		log.Printf("Failed to save notes: %v", err)
	}
}

func Delete(id int) {
	for i, n := range notes {
		if n.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			if err := store.Save(notes); err != nil {
				log.Printf("Failed to save notes: %v", err)
			}
			return
		}
	}
}

func Update(id int, content string) {
	for i := range notes {
		if notes[i].ID == id {
			// Found the note, update it in place.
			notes[i].Content = content
			notes[i].Created = time.Now()

			// Regenerate the embedding for the new content.
			embedding, err := ai.GenerateEmbedding(content)
			if err != nil {
				log.Printf("Failed to generate embedding during update: %v", err)
			}
			notes[i].Embedding = embedding

			// Save the entire updated slice back to the store.
			if err := store.Save(notes); err != nil {
				log.Printf("Failed to save notes after update: %v", err)
			}
			return // Exit the function
		}
	}
}

func All() []Note {
	return notes
}
