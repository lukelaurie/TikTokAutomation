package uploadVideoHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
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

	audioPath := api.GenerateAudioFile()

	generateTikTokVideo(videoText, audioPath, videoPath, fontName, fontColor)

	json.NewEncoder(w).Encode("yooooooo")
}

func generateTikTokVideo(videoText string, audioPath string, videoPath string, fontName string, fontColor string) {
	// specify the path for the output file
	outputPath := "./assets/video/output.mp4"
	cutPath := "./assets/video/cut.mp4"
	combinedPath := "./assets/video/combined.mp4"

	// have video and mp3 length match each other
	cutErr := cutVideoLength(audioPath, videoPath, cutPath)
	if cutErr != nil {
		panic(cutErr)
	}

	// combine the mp4 and audio files into one video
	combineErr := combineAudioAndVideo(audioPath, cutPath, combinedPath)
	if combineErr != nil {
		panic(combineErr)
	}

	// add the text from the api onto the video
	overlayErr := overlayTextOnVideo(videoText, fontName, fontColor, combinedPath, outputPath)
	if overlayErr != nil {
		panic(overlayErr)
	}

	fmt.Println("combined")
}

func cutVideoLength(audioPath string, videoPath string, cutPath string) error {
	cmdGetDuration := exec.Command("ffmpeg.exe", "-i", audioPath, "-f", "null", "-")
	durationOutput, err := cmdGetDuration.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg error: %v\nOutput: %s", err, durationOutput)
	}

	// Extract the duration from the output
	re := regexp.MustCompile(`Duration: (\d+:\d+:\d+\.\d+)`)
	match := re.FindStringSubmatch(string(durationOutput))
	if len(match) < 2 {
		return fmt.Errorf("failed to parse audio from mp3 output")
	}
	duration := match[1]

	// cut down the length of the video
	cmdCutVideo := exec.Command("ffmpeg.exe", "-i", videoPath, "-t", duration, "-c", "copy", "-y", cutPath)
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

func overlayTextOnVideo(videoText string, fontName string, fontColor string, combinedPath string, outputPath string) error {
	textIntervals := generateTextDisplayIntervals(videoText)

	// filter the string to place all start and end dates
	var filters []string
	fontFile := fmt.Sprintf("./assets/font/%s.ttf", fontName)

	for _, textInterval := range textIntervals {
		textFormat := fmt.Sprintf("drawtext=text='%s':fontfile=%s:fontsize=60:fontcolor=%s:x=(w-text_w)/2:y=(h-text_h)-700:enable='between(t,%f,%f)'",
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

func generateTextDisplayIntervals(videoText string) []model.TextDisplay {
	var allText []model.TextDisplay
	allWords := strings.Split(videoText, " ")
	var curTime float32 = 0.5

	// generate the strings to place the text on the video
	for i := 0; i < len(allWords); i++ {
		curWord := allWords[i]
		// determine weather to show 1 or 2 words at a time
		if len(curWord) < 5 && i < len(allWords) - 1 {
			curWord += " " + allWords[i + 1]
			i++;
		} 
		const timeAdd float32 = 0.6
		allText = append(allText, model.TextDisplay{Text: curWord, StartTime: curTime, EndTime: curTime + timeAdd})
		curTime += timeAdd
	}

	return allText
}