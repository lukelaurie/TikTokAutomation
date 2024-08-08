package model 

type VideoPromptData struct {
    Instructions string   `json:"instructions"`
    Prompts      []string `json:"prompts"`
}

type VideoPrompts map[string]VideoPromptData