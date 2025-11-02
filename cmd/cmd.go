package cmd

import (
	"flag"
	"fmt"
	"os"

	"note-app/ai"
	"note-app/note"
	"note-app/store"
)

func Execute() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addContent := addCmd.String("content", "", "The content of the note")

	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addContent == "" {
			fmt.Println("Please provide content for the note")
			os.Exit(1)
		}
		note.Add(*addContent)
		fmt.Println("Note added successfully")
		os.Exit(0)
	case "regenerate":
		notes := note.All()
		for i, n := range notes {
			embedding, err := ai.GenerateEmbedding(n.Content)
			if err != nil {
				fmt.Printf("Failed to generate embedding for note %d: %v\n", n.ID, err)
				continue
			}
			notes[i].Embedding = embedding
		}
		if err := store.Save(notes); err != nil {
			fmt.Printf("Failed to save notes: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Embeddings regenerated successfully")
		os.Exit(0)
	}
}
