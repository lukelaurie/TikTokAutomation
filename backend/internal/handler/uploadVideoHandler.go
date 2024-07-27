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
	videoText := api.GenerateVideoScript(videoType)
	videoText += ""
	audioError := api.GenerateAudioFile(videoText)
	if audioError != nil {
		panic(audioError)
	}
	audioPath := "./assets/audio/output.wav"

	videoLength, allText, timeStampError := api.TimeStampGenerator(audioPath)

	if timeStampError != nil {
		panic(timeStampError)
	}

	generateTikTokVideo(audioPath, videoPath, fontName, fontColor, videoLength, allText)

	json.NewEncoder(w).Encode("video created")
}

func generateTikTokVideo(audioPath string, videoPath string, fontName string, fontColor string, videoLength float64, allText *[]model.TextDisplay) {
	// specify the path for the output file
	outputPath := "./assets/video/output.mp4"
	cutPath := "./assets/video/cut.mp4"
	combinedPath := "./assets/video/combined.mp4"

	// have video and mp3 length match each other
	cutErr := cutVideoLength(videoPath, cutPath, videoLength)
	if cutErr != nil {
		panic(cutErr)
	}

	// combine the mp4 and audio files into one video
	combineErr := combineAudioAndVideo(audioPath, cutPath, combinedPath)
	if combineErr != nil {
		panic(combineErr)
	}

	// add the text from the api onto the video
	overlayErr := overlayTextOnVideo(fontName, fontColor, combinedPath, outputPath, allText)
	if overlayErr != nil {
		panic(overlayErr)
	}
	removeFileErr := deleteFile(cutPath, combinedPath)
	if removeFileErr != nil {
		panic(removeFileErr)
	}
}

func cutVideoLength(videoPath string, cutPath string, videoLength float64) error {
	// cut down the length of the video
	strVideoLength := fmt.Sprintf("%f", videoLength)
	cmdCutVideo := exec.Command("ffmpeg.exe", "-i", videoPath, "-t", strVideoLength, "-c", "copy", "-y", cutPath)
	cutOutput, err := cmdCutVideo.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, cutOutput)
	}

	return nil
}

func combineAudioAndVideo(audioPath string, videoPath string, outputPath string) error {
	cmd := exec.Command("ffmpeg.exe", "-i", videoPath, "-i", audioPath, "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "-y", outputPath)

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
		textFormat := fmt.Sprintf("drawtext=text='%s':fontfile=%s:fontsize=60:fontcolor=%s:x=(w-text_w)/2:y=(h-text_h)-700:enable='between(t,%s,%s)'",
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

func deleteFile(cutPath string, combinedPath string) error {
	cutErr := os.Remove(cutPath)
	combineErr := os.Remove(combinedPath)
	if cutErr != nil || combineErr != nil {
		return fmt.Errorf("error removing file")
	}
	return nil
}