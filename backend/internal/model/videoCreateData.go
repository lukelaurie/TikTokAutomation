package model 

// info for the text that gets displayed on the video
type TextDisplay struct {
	Text 	  string
	StartTime string
	EndTime   string
}

type VideoCreationInfo struct {
	AudioPath           string
	VideoPath           string
	BackgroundAudioPath string
	AllText             *[]TextDisplay
}

// what the json file with all the prompts gets converted into
type VideoPromptData struct {
    Instructions string   `json:"instructions"`
    Prompts      []string `json:"prompts"`
}

type VideoPrompts map[string]VideoPromptData