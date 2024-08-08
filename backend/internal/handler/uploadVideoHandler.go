package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"

	api "github.com/lukelaurie/TikTokAutomation/backend/internal/api"
	database "github.com/lukelaurie/TikTokAutomation/backend/internal/database"
	model "github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	// TODO -> Get all user with passed in membership type

	// TODO -> Determine which preference to use

	// retrieve the data from the needed scheduler
	videoType, backgroundType, fontName, fontColor := database.RetrieveSchedulerInfo(5)
	videoPath, videoPathErr := getRandomBackgroundFile(backgroundType, "video")
	if videoPathErr != nil {
		panic(videoPathErr)
	}

	backgroundAudioPath, videoPathErr := getRandomBackgroundFile(videoType, "audio")
	if videoPathErr != nil {
		panic(videoPathErr)
	}

	backgroundAudioPath += ""

	videoText, scriptErr := api.GenerateVideoScript(videoType)
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

	generateTikTokVideo(audioPath, videoPath, backgroundAudioPath, fontName, fontColor, allText)

	json.NewEncoder(w).Encode("success")
}

func getRandomBackgroundFile(backgroundVideoType string, backgroundType string) (string, error) {
	directoryPath := fmt.Sprintf("./assets/%s/%s/", backgroundType, backgroundVideoType)

	// read the directory
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return "", fmt.Errorf("error opening directoy: %v", err)
	}

	// collect all of the file names
	var fileNames []string
	for _, file := range files {
		if file.Type().IsRegular() {
			fileNames = append(fileNames, file.Name())
		}
	}

	// verify files existed in the directory
	if len(fileNames) == 0 {
		return "", fmt.Errorf("error: no files exist in the firectory")
	}

	randIndex := rand.Intn(len(fileNames))
	fileName := fileNames[randIndex]

	return directoryPath + fileName, nil
}

func generateTikTokVideo(audioPath string, videoPath string, backgroundAudioPath string, fontName string, fontColor string, allText *[]model.TextDisplay) {
	// specify the path for the used files in the combination
	outputPath := "./assets/video/output.mp4"
	combinedPath := "./assets/video/combined.mp4"
	combinedAudioPath := "./assets/audio/combined-output.wav"

	// combine the two audio files
	audioCombineErr := combineAudioFiles(audioPath, backgroundAudioPath, combinedAudioPath)
	if audioCombineErr != nil {
		panic(audioCombineErr)
	}

	// combine the mp4 and audio files into one video
	combineErr := combineAudioAndVideo(combinedAudioPath, videoPath, combinedPath)
	if combineErr != nil {
		panic(combineErr)
	}

	// add the text from the api onto the video
	overlayErr := overlayTextOnVideo(fontName, fontColor, combinedPath, outputPath, allText)
	if overlayErr != nil {
		panic(overlayErr)
	}
	// delete the audio and video files
	removeErr := deleteFiles(combinedPath, audioPath, combinedAudioPath)
	if removeErr != nil {
		panic(removeErr)
	}
}

func combineAudioFiles(audioPath string, backgroundAudioPath string, outputPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", audioPath,
		"-i", backgroundAudioPath,
		"-filter_complex", "[0:a][1:a]amix=inputs=2:duration=longest[a]",
		"-map", "[a]",
		"-c:a", "pcm_s16le", // Set the audio codec to PCM for WAV format
		"-y", outputPath,
	)

	// Run the command and get the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg error combining audio files: %v\nOutput: %s", err, output)
	}

	return nil
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
		textFormat := fmt.Sprintf(
			"drawtext=text='%s':fontfile=%s:fontsize=w/6:fontcolor=%s:x=(w-text_w)/2:y=(h-text_h)/2:shadowx=2:shadowy=2:shadowcolor=black:enable='between(t,%s,%s)'",
			textInterval.Text, fontFile, fontColor, textInterval.StartTime, textInterval.EndTime)

		filters = append(filters, textFormat)
	}

	filteredString := strings.Join(filters, ",")
	cmdOverlayText := exec.Command("ffmpeg.exe", "-i", combinedPath, "-vf", filteredString, "-c:a", "copy", "-y", outputPath)
	overlayOutput, err := cmdOverlayText.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, overlayOutput)
	}

	return nil
}

func deleteFiles(combinedPath string, audioPath string, combinedAudioPath string) error {
	combineErr := os.Remove(combinedPath)
	if combineErr != nil {
		return fmt.Errorf("error removing file")
	}
	outputErr := os.Remove(audioPath)
	if outputErr != nil {
		return fmt.Errorf("error removing file")
	}
	combineAudioErr := os.Remove(combinedAudioPath)
	if combineAudioErr != nil {
		return fmt.Errorf("error removing file")
	}
	return nil
}
