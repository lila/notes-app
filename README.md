# Note App

Terminal note taking app with concept grouping using Google Gemini.

## Features
- Journal view by creation date
- Concept view groups notes by embedding similarity
- Live editing with regenerating embeddings

## Prerequisites
- Go 1.22+
- `GEMINI_API_KEY` environment variable set for Gemini SDK

## Usage
```bash
# Run the TUI
go run .

# Add a note from the CLI
go run . add --content "Your note"

# Regenerate embeddings for all notes
go run . regenerate
```

## Concept Mode Refresh
Hit `Tab` to switch to concepts, then press `r` to clear the current list, show
`Regenerating concept groupings...`, and invoke Gemini for fresh names (cache is
ignored during this forced refresh).

## Storage
- Notes persist in `notes.json`
- Concept names cache in `concept_cache.json`
