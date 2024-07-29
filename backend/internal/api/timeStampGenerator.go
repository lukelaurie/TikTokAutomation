package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	model "github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func TimeStampGenerator(audioPath string) (float64, *[]model.TextDisplay, error) {
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	// endpoint := fmt.Sprintf("https://%s.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=en-US", speechRegion)
	endpoint := fmt.Sprintf("https://%s.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=en-US&format=detailed&profanity=masked&wordLevelTimestamps=true", speechRegion)


	audioFile, err := os.ReadFile(audioPath)
	if err != nil {
		return -1, nil, fmt.Errorf("error getting the audio file: %v", err)
	}

	// generate the request
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(audioFile))

	if err != nil {
		return -1, nil, fmt.Errorf("azure error making the request: %v", err)
	}

	req.Header.Set("Content-Type", "audio/wav")
	req.Header.Set("Ocp-Apim-Subscription-Key", speechKey)

	// execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, fmt.Errorf("azure error sending the request: %v", err)
	}
	// close the response after function call completes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		return -1, nil, fmt.Errorf("response Body: %s", body)
	}

	// Unmarshal the response JSON into the structured data
	var transcriptResponse model.TranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&transcriptResponse); err != nil {
		return -1, nil, fmt.Errorf("error decoding azure response body: %v", err)
	}

	videoLength, allText := processTranscriptResponse(&transcriptResponse)

	return videoLength, allText, nil
}

func convertNano(time int) float64 {
	return float64(time) / 10000000
}
func processTranscriptResponse(transcriptResponse *model.TranscriptResponse) (float64, *[]model.TextDisplay) {
	videoLength := convertNano(transcriptResponse.Duration)
	// extract all words and their time stamps 
	var allText []model.TextDisplay
	allWords := transcriptResponse.NBest[0].Words
	for i := 0; i < len(allWords); i++ {
		wordItem := allWords[i]
		// parse the text and start/end times
		curWord := wordItem.Word
		startTime := convertNano(wordItem.Offset)
		duration := convertNano(wordItem.Duration)
		endTime := startTime + duration
		// check if need to show 1 or 2 words
		if len(curWord) < 5 && i < len(allWords) - 1 {
			curWord += " " + allWords[i + 1].Word
			newDuration := convertNano(allWords[i + 1].Duration)
			endTime += newDuration
			i++;
		}
		// ` causes string to be invalid 
		curWord = strings.ReplaceAll(curWord, "'", "\\")
		// have word come up a little before spoken 
		startTime -= 0.25
		endTime -= 0.25
		allText = append(allText, model.TextDisplay{
			Text: curWord, 
			StartTime: fmt.Sprintf("%.1f", startTime), 
			EndTime: fmt.Sprintf("%.1f", endTime)}) 

	}
	return videoLength, &allText
}