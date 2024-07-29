package model

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