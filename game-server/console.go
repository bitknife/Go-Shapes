package main

import (
	"fmt"
	// https://github.com/c-bata/go-prompt
	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "shutdown()", Description: "Shutdown gracefully."},
		{Text: "quit()", Description: "Quit without graceful shutdown."},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func StartConsole() {
	for {
		fmt.Println("Welcome to the server, hit TAB to see commands.")
		t := prompt.Input("> ", completer)
		fmt.Println("You selected " + t)
	}
}
