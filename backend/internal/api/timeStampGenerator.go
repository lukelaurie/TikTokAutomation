package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)
func TimeStampGenerator(audioPath string) (*[]model.TextDisplay, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	audioFile, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("error opening the audio file: %v", err)
	}	
	defer audioFile.Close()

	var requstBody bytes.Buffer
	writer := multipart.NewWriter(&requstBody)

	// create form in the body and add the audio file to it
	part, err := writer.CreateFormFile("file", audioPath)
	if err != nil {
		return nil, fmt.Errorf("error creating the audio form: %v", err)
	}	
	_, err = io.Copy(part, audioFile)
	if err != nil {
		return nil, fmt.Errorf("error copying the audio file: %v", err)
	}	

	// write the reamaining parts of the body
	writer.WriteField("model", "whisper-1")
	writer.WriteField("response_format", "verbose_json")
	writer.WriteField("timestamp_granularities[]", "word")

	writer.Close()
	
	// generate the request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &requstBody)
	if err != nil {
		return nil, fmt.Errorf("azure error making the request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer " + apiKey)

	// execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("azure error sending the request: %v", err)
	}
	// close the response after function call completes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("response Body: %s", body)
	}

	// Unmarshal the response JSON into the structured data
	var transcriptResponse model.TranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&transcriptResponse); err != nil {
		return nil, fmt.Errorf("error decoding azure response body: %v", err)
	}

	allText := processTranscriptResponse(&transcriptResponse)

	return allText, nil
}

func processTranscriptResponse(transcriptResponse *model.TranscriptResponse) *[]model.TextDisplay {
	// extract all words and their time stamps
	var allText []model.TextDisplay
	allWords := transcriptResponse.Words
	for i := 0; i < len(allWords); i++ {
		wordItem := allWords[i]
		// parse the text and start/end times
		curWord, startTime, endTime := wordItem.Word, wordItem.Start, wordItem.End
		if len(curWord) < 5 && i < len(allWords)-1 {
			curWord += " " + allWords[i+1].Word
			endTime = allWords[i+1].End
			i++
		}
		
		// ` causes string to be invalid
		curWord = strings.ReplaceAll(curWord, "'", "\\")

		// have word come up a little before spoken
		startTime -= 0.25
		endTime -= 0.25
		startTime = math.Max(0, startTime)
		endTime = math.Max(0, endTime)

		allText = append(allText, model.TextDisplay{
			Text:      curWord,
			StartTime: fmt.Sprintf("%.1f", startTime),
			EndTime:   fmt.Sprintf("%.1f", endTime)})
	}
	return &allText
}