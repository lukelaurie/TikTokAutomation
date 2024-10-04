package model

// structs to interact with chatGPT
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatGPTRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type Choice struct {
    Message Message `json:"message"`
}

type ChatGPTResponse struct {
    Choices []Choice `json:"choices"`
}

// structs to interact with azure to turn speech into text
type TranscriptResponse struct {
	Task      string    `json:"task"`
	Language  string    `json:"language"`
	Duration  float64   `json:"duration"`
	Text      string    `json:"text"`
	Words     []Word    `json:"words"`
}

type Word struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}