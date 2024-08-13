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
	"github.com/lukelaurie/TikTokAutomation/backend/internal/middleware"
	model "github.com/lukelaurie/TikTokAutomation/backend/internal/model"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/utils"
)

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	// TODO -> Get all user with passed in membership type

	// this is TEMPORARY to get user from cookie. In future get all for membership type. And execute
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "Username not found in context", http.StatusInternalServerError)
		return
	}

	preference, dbErr := getPreferenceFromDatabase(username)
	if dbErr != nil {
		utils.LogAndAddServerError(dbErr, w)
        return
	}

	videoCreationInfo, videoErr := getVideoCreationInformation(preference)
	if videoErr != nil {
		utils.LogAndAddServerError(videoErr, w)
        return
	}


	videoCombineError := generateTikTokVideo(preference, videoCreationInfo)
	if videoCombineError != nil {
		utils.LogAndAddServerError(videoCombineError, w)
        return
	}

	json.NewEncoder(w).Encode("success")
}

func getPreferenceFromDatabase(username string) (model.Preference, error) {
	// determine which preference to use for the user
	preferenceTrackerIndex, preferenceIndexErr := retrieveAndUpdatePreferenceIndex(username)
	if preferenceIndexErr != nil {
		return model.Preference{}, preferenceIndexErr
	}

	// retrieve the data from the needed scheduler
	preference, preferenceErr := database.RetrieveSchedulerInfo(username, preferenceTrackerIndex)
	if preferenceErr != nil {
		return model.Preference{}, preferenceIndexErr
	}

	return preference, nil
}

func getVideoCreationInformation(preference model.Preference) (model.VideoCreationInfo, error) {
	var videoCreationInfo model.VideoCreationInfo

	videoPath, videoPathErr := getRandomBackgroundFile(preference.BackgroundVideoType, "video")
	if videoPathErr != nil {
        return model.VideoCreationInfo{}, videoPathErr
	}
	videoCreationInfo.VideoPath = videoPath

	backgroundAudioPath, videoPathErr := getRandomBackgroundFile(preference.VideoType, "audio")
	if videoPathErr != nil {
		return model.VideoCreationInfo{}, videoPathErr
	}
	videoCreationInfo.BackgroundAudioPath = backgroundAudioPath

	videoText, scriptErr := api.GenerateVideoScript(preference.VideoType)
	if scriptErr != nil {
		return model.VideoCreationInfo{}, scriptErr
	}

	audioError := api.GenerateAudioFile(videoText)
	if audioError != nil {
		return model.VideoCreationInfo{}, audioError
	}

	audioPath := "./assets/audio/output.wav"
	allText, timeStampError := api.TimeStampGenerator(audioPath)
	if timeStampError != nil {
		return model.VideoCreationInfo{}, timeStampError
	}
	videoCreationInfo.AudioPath = audioPath
	videoCreationInfo.AllText = allText

	return videoCreationInfo, nil
}

func retrieveAndUpdatePreferenceIndex(username string) (int, error) {
	// get the current index for the preference
	preferenceTracker, err := database.RetrievePreferenceTracker(username)
	if err != nil {
		return -1, err
	}

	if preferenceTracker.CurPreferenceCount == 0 {
		return -1, fmt.Errorf("user: %s has not yet created a preference", username)
	}

	// increment the preference to next index or restart to front if already created a video for all 
	newPreferenceIndex := preferenceTracker.CurPreferenceOrder + 1
	if newPreferenceIndex > preferenceTracker.CurPreferenceCount {
		newPreferenceIndex = 1
	}
	
	// increment the count in the database by one 
	err = database.IncrementPreferenceTracker(username, newPreferenceIndex, false)
	if err != nil {
		return -1, fmt.Errorf("error incrementing preference index in the database")
	}

	return preferenceTracker.CurPreferenceOrder, nil
}

func getRandomBackgroundFile(backgroundVideoType string, backgroundType string) (string, error) {
	directoryPath := fmt.Sprintf("./assets/%s/%s/", backgroundType, backgroundVideoType)

	// read the directory
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return "", fmt.Errorf("error opening directory: %v", err)
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
		return "", fmt.Errorf("error: no files exist in the directory")
	}

	randIndex := rand.Intn(len(fileNames))
	fileName := fileNames[randIndex]

	return directoryPath + fileName, nil
}

func generateTikTokVideo(preference model.Preference, videoCreationInfo model.VideoCreationInfo) error {
	// specify the path for the used files in the combination
	outputPath := "./assets/video/output.mp4"
	combinedPath := "./assets/video/combined.mp4"
	combinedAudioPath := "./assets/audio/combined-output.wav"

	// combine the two audio files
	audioCombineErr := combineAudioFiles(videoCreationInfo.AudioPath, videoCreationInfo.BackgroundAudioPath, combinedAudioPath)
	if audioCombineErr != nil {
		return audioCombineErr
	}

	// combine the mp4 and audio files into one video
	combineErr := combineAudioAndVideo(combinedAudioPath, videoCreationInfo.VideoPath, combinedPath)
	if combineErr != nil {
		return combineErr
	}

	// add the text from the api onto the video
	overlayErr := overlayTextOnVideo(preference.FontName, preference.FontColor, combinedPath, outputPath, videoCreationInfo.AllText)
	if overlayErr != nil {
		return overlayErr
	}
	// delete the audio and video files
	removeErr := deleteFiles(combinedPath, videoCreationInfo.AudioPath, combinedAudioPath)
	if removeErr != nil {
		return removeErr
	}
	return nil
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
