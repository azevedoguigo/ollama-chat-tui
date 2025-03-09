package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/azevedoguigo/ollama-chat-tui/internal/storage"
)

type ModelDetails struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

type Model struct {
	Name       string       `json:"name"`
	Model      string       `json:"model"`
	ModifiedAt time.Time    `json:"modified_at"`
	Size       int64        `json:"size"`
	Digest     string       `json:"digest"`
	Details    ModelDetails `json:"details"`
}

type FindModelsResponse struct {
	Models []Model `json:"models"`
}

type OllamaRequest struct {
	Model    string            `json:"model"`
	Prompt   string            `json:"prompt"`
	Stream   bool              `json:"stream"`
	Messages []storage.Message `json:"messages"`
}

func FindOllamaLocalModels() (*FindModelsResponse, error) {
	response, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var findModelsResponse FindModelsResponse

	err = json.Unmarshal(data, &findModelsResponse)
	if err != nil {
		return nil, err
	}

	return &findModelsResponse, nil
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
