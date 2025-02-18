package main

import (
	"log"

	"github.com/azevedoguigo/ollama-chat-tui/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Erro ao rodar a aplicação: %v", err)
	}
}
