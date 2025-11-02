package main

import (
	"fmt"
	"os"

	"note-app/cmd"
	"note-app/tui"
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	p := tui.NewProgram()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
