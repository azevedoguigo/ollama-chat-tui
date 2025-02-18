package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/azevedoguigo/ollama-chat-tui/internal/storage"
)

type OllamaRequest struct {
	Model    string            `json:"model"`
	Prompt   string            `json:"prompt"`
	Stream   bool              `json:"stream"`
	Messages []storage.Message `json:"messages"`
}

func QueryOllamaStream(model string, messages []storage.Message, callbalc func(string)) error {
	requestData := OllamaRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	response, err := http.Post(
		"http://localhost:11434/api/chat",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		var data map[string]interface{}

		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			return err
		}

		if message, ok := data["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				callbalc(content)
			}
		}
	}

	return scanner.Err()
}
