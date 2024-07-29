package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func GenerateVideoScript(videoType string, chatInstructions string, chatPrompt string) (string, error) {
	const apiUrl = "https://api.openai.com/v1/chat/completions"

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("error: please provide the enironment variable \"OPENAI_API_KEY\"")
	}

	// provide instructions on how the gpt model should respond
	messages := []model.Message{
		{Role: "system", Content: chatInstructions},
		{Role: "user", Content: chatPrompt},
	}

	reqeustBody := model.ChatGPTRequest{
		Model:    "gpt-4o-mini",
		Messages: messages,
	}

	// convert message to json befoer sending to gpt
	jsonRequest, err := json.Marshal(reqeustBody)
	if err != nil {
		return "", fmt.Errorf("error converting the request to json: %v", err)
	}

	// generate the request
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return "", fmt.Errorf("error narshaling the request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("chatGPT error sending the request: %v", err)
	}
	// close the response after function call completes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("response Body: %s", body)
	}

	var chatResponse model.ChatGPTResponse 
	err = json.NewDecoder(resp.Body).Decode(&chatResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}
	if len(chatResponse.Choices) <= 0 {
		return "", fmt.Errorf("no response from chatGPT")
	}
	return chatResponse.Choices[0].Message.Content, nil
}