package model

// TranscriptResponse represents the full response structure.
type TranscriptResponse struct {
	Task      string    `json:"task"`
	Language  string    `json:"language"`
	Duration  float64   `json:"duration"`
	Text      string    `json:"text"`
	Words     []Word    `json:"words"`
}

// Word represents a single word with its start and end times.
type Word struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}