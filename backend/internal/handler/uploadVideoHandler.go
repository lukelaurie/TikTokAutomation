package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	// "regexp"
	"strings"

	api "github.com/lukelaurie/TikTokAutomation/backend/internal/api"
	database "github.com/lukelaurie/TikTokAutomation/backend/internal/database"
	model "github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	// TODO -> Get all user with passed in membership type

	// TODO -> Determine which preference to use

	// retrieve the data from the needed scheduler
	videoType, videoPath, fontName, fontColor := database.RetrieveSchedulerInfo(5)


	chatInstructions := "You are a compelling story teller. Each story you generate must be unique from any other story told, and should be around 30 seconds long"
	chatPrompt := "Connor is a avid drug addict. Tell me a story about Connor going to the gym and tearing his bicep terribly, then dying and people laughing at him through the whole process."
	videoText, scriptErr := api.GenerateVideoScript(videoType, chatInstructions, chatPrompt)
	if scriptErr != nil {
		panic(scriptErr)
	}
	audioError := api.GenerateAudioFile(videoText)
	if audioError != nil {
		panic(audioError)
	}
	audioPath := "./assets/audio/output.wav"

	allText, timeStampError := api.TimeStampGenerator(audioPath)

	if timeStampError != nil {
		panic(timeStampError)
	}

	generateTikTokVideo(audioPath, videoPath, fontName, fontColor, allText)

	json.NewEncoder(w).Encode("success")
}

func generateTikTokVideo(audioPath string, videoPath string, fontName string, fontColor string, allText *[]model.TextDisplay) {
	// specify the path for the output file
	outputPath := "./assets/video/output.mp4"
	combinedPath := "./assets/video/combined.mp4"

	// combine the mp4 and audio files into one video
	combineErr := combineAudioAndVideo(audioPath, videoPath, combinedPath)
	if combineErr != nil {
		panic(combineErr)
	}

	// add the text from the api onto the video
	overlayErr := overlayTextOnVideo(fontName, fontColor, combinedPath, outputPath, allText)
	if overlayErr != nil {
		panic(overlayErr)
	}
	removeFileErr := deleteFile(combinedPath)
	if removeFileErr != nil {
		panic(removeFileErr)
	}
}

func combineAudioAndVideo(audioPath string, videoPath string, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "-shortest", "-y", outputPath)

	// Run the command and get the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg error: %v\nOutput: %s", err, output)
	}

	return nil
}

func overlayTextOnVideo(fontName string, fontColor string, combinedPath string, outputPath string, allText *[]model.TextDisplay) error {
	// textIntervals := generateTextDisplayIntervals(videoText)
	textIntervals := *allText

	// filter the string to place all start and end dates
	var filters []string
	fontFile := fmt.Sprintf("./assets/font/%s.ttf", fontName)

	for _, textInterval := range textIntervals {
		textFormat := fmt.Sprintf("drawtext=text='%s':fontfile=%s:fontsize=w/7:fontcolor=%s:x=(w-text_w)/2:y=(h-text_h)/3:enable='between(t,%s,%s)'",
		textInterval.Text, fontFile, fontColor, textInterval.StartTime, textInterval.EndTime)

		filters = append(filters, textFormat)
	}	

	filteredString := strings.Join(filters, ",")
	// fmt.Printf("%s", filteredString)
	cmdOverlayText := exec.Command("ffmpeg.exe", "-i", combinedPath, "-vf", filteredString, "-c:a", "copy", "-y", outputPath)
	overlayOutput, err := cmdOverlayText.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, overlayOutput)
	}

	// delete the remaining files 

	return nil
}

func deleteFile(combinedPath string) error {
	combineErr := os.Remove(combinedPath)
	if combineErr != nil {
		return fmt.Errorf("error removing file")
	}
	return nil
}