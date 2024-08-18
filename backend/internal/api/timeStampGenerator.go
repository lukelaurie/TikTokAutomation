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

	"github.com/agnivade/levenshtein"

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

	allWordsIndex := 0
	splitTextIndex := 0
	// for i := 0; i < len(allWords) && i < len(splitText); i++ {
	for allWordsIndex < len(allWords) && splitTextIndex < len(splitText) {
		if !checkIfWordsSame(allWords[allWordsIndex].Word, splitText[splitTextIndex]) {
			// check if splitText was missing a word
			if splitTextIndex+1 < len(splitText) && checkIfWordsSame(allWords[allWordsIndex].Word, splitText[splitTextIndex+1]) {
				splitTextIndex += 1
			}
			// check if allWords was missing a word
			if allWordsIndex+1 < len(allWords) && checkIfWordsSame(allWords[allWordsIndex+1].Word, splitText[splitTextIndex]) {
				allWordsIndex += 1
			}
		}

		curWord := splitText[splitTextIndex]
		// add new words while they fit into the screen and not end of sentence
		for !strings.HasSuffix(curWord, ".") &&
			splitTextIndex < len(splitText)-1 &&
			(len(curWord)+len(splitText[splitTextIndex+1])+1) <= 11 {
			curWord += " " + splitText[splitTextIndex+1]
			splitTextIndex += 1
			allWordsIndex += 1
		}

		// use the previous end time of prior word to account for whisper miscalculations
		startTime := endTime
		endTime = allWords[allWordsIndex].End

		// add to end time if end of sentence to account for the pause
		if strings.HasSuffix(curWord, ".") {
			endTime += .7
		}

		// ' causes string to be invalid
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
		
		splitTextIndex += 1
		allWordsIndex += 1
	}
	return &allText
}

func checkIfWordsSame(word1 string, word2 string) bool {
	// Calculate the Levenshtein distance between the two words
	distance := levenshtein.ComputeDistance(word1, word2)
	maxLen := max(len(word1), len(word2))

	// Calculate the similarity percentage
	similarity := (1 - float64(distance)/float64(maxLen)) * 100
	return similarity > 50
}
