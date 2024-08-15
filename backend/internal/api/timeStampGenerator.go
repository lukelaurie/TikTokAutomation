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
	req.Header.Set("Authorization", "Bearer "+apiKey)

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
	// use the split text over the words because it contains grammar
	splitText := strings.Split(transcriptResponse.Text, " ")

	// extract all words and their time stamps
	var allText []model.TextDisplay
	allWords := transcriptResponse.Words
	endTime := 0.0

	for i := 0; i < len(allWords) && i < len(splitText); i++ {
		// use the previous end time of prior word to account for whisper miscalculations
		startTime := endTime
		curWord := splitText[i]

		// add new words while they fit into the screen and not end of sentence
		for !strings.HasSuffix(curWord, ".") &&
			i < len(allWords)-1 &&
			(len(curWord)+len(splitText[i+1])+1) <= 11 {
			curWord += " " + splitText[i+1]
			i++
		}
		endTime = allWords[i].End
		// add to end time if end of sentence to account for the pause
		if strings.HasSuffix(curWord, ".") {
			endTime += .7
		}

		// ` causes string to be invalid
		curWord = strings.ReplaceAll(curWord, "'", "’")

		// have word come up a little before spoken
		earlyStartTime := math.Max(0, startTime-0.3)
		earlyEndTime := math.Max(0, endTime-0.3)

		allText = append(allText, model.TextDisplay{
			Text:      curWord,
			StartTime: fmt.Sprintf("%.2f", earlyStartTime),
			EndTime:   fmt.Sprintf("%.2f", earlyEndTime)})
		// this generates the text to be placed in test text
		// fmt.Printf("%s:%s:%s^^", curWord, fmt.Sprintf("%.2f", earlyStartTime), fmt.Sprintf("%.2f", earlyEndTime))
	}
	return &allText
}
