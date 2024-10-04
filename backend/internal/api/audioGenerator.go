package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GenerateAudioFile(videoText string, voice string) error {
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	endpoint := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", speechRegion)

	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xml:lang="en-US">
		<voice xml:lang="en-US" name="%s">
			%s
		</voice>
	</speak>`, voice, videoText)


	// generate the request
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(body))

	if err != nil {
		return fmt.Errorf("azure error making the request: %v", err)
	}

	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", speechKey)
    req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")

	// execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("azure error sending the request: %v", err)
	}
	// close the response after function call completes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("response Body: %s", body)
	}

	// convert the response into an audio file 
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading azure response body: %v", err)
	}

	err = os.WriteFile("./assets/audio/output.wav", audioData, 0644) // 0644 sets file permisisons
	if err != nil {
		return fmt.Errorf("error converting body to .wav file: %v", err)
	}
	return nil
}