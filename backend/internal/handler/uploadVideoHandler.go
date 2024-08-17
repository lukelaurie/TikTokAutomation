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

func UploadVideo(isTestMode bool, w http.ResponseWriter, r *http.Request) {
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

	videoCreationInfo, videoErr := getVideoCreationInformation(preference, isTestMode)
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

func getVideoCreationInformation(preference model.Preference, isTestMode bool) (model.VideoCreationInfo, error) {
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

	var audioPath string
	var allText *[]model.TextDisplay

	if isTestMode {
		// don't call api's and just use test data
		textAndTimes := "In a small:0.00:0.34^^town, Fiona:0.34:1.28^^was:1.28:1.54^^abducted by:1.54:2.12^^a man who:2.12:2.76^^lived in a:2.76:3.22^^decrepit:3.22:3.56^^house on:3.56:4.04^^the:4.04:4.24^^outskirts.:4.24:5.40^^For ten:5.40:6.28^^long years,:6.28:6.90^^she was:6.90:7.48^^trapped:7.48:7.74^^within its:7.74:8.30^^dark walls,:8.30:8.96^^a prisoner:8.96:9.60^^to his:9.60:10.02^^twisted:10.02:10.32^^games.:10.32:11.50^^Each day,:11.50:12.40^^she made:12.40:12.96^^plans to:12.96:13.60^^escape,:13.60:13.94^^plotting:13.94:14.46^^intricate:14.46:14.94^^routes:14.94:15.22^^through the:15.22:15.72^^house when:15.72:16.12^^her captor:16.12:16.60^^was absent.:16.60:17.96^^One dreary:17.96:18.88^^evening,:18.88:19.28^^after years:19.28:20.14^^of failed:20.14:20.56^^attempts,:20.56:21.04^^she finally:21.04:21.98^^found a:21.98:22.62^^window:22.62:22.62^^unlocked.:22.62:23.86^^Heart:23.86:24.46^^racing, she:24.46:25.28^^pushed it:25.28:25.70^^open, the:25.70:26.38^^cold night:26.38:26.86^^air,:26.86:27.16^^beckoning:27.16:27.60^^her towards:27.60:28.18^^freedom.:28.18:29.26^^As she:29.26:29.92^^crept:29.92:30.18^^outside,:30.18:30.70^^her feet:30.70:31.30^^barely:31.30:31.68^^touching:31.68:32.00^^the ground,:32.00:32.58^^the sound:32.58:33.16^^of:33.16:33.34^^footsteps:33.34:33.78^^echoed in:33.78:34.26^^the quiet.:34.26:35.44^^Before she:35.44:36.32^^could run,:36.32:36.82^^she felt a:36.82:37.76^^hand grasp:37.76:38.12^^her:38.12:38.38^^shoulder.:38.38:39.36^^You never:39.36:40.14^^learn, do:40.14:40.88^^you? he:40.88:41.50^^whispered:41.50:41.84^^in her ear,:41.84:42.38^^a sinister:42.38:43.16^^smile:43.16:43.62^^spreading:43.62:44.00^^across his:44.00:44.60^^face.:44.60:45.74^^Overwhelmed:45.74:46.54^^with:46.54:46.78^^despair,:46.78:47.22^^Fiona:47.22:47.90^^realized:47.90:48.46^^that her:48.46:48.90^^captor had:48.90:49.42^^been:49.42:49.60^^watching:49.60:49.90^^her every:49.90:50.42^^move,:50.42:50.76^^delighting:50.76:51.54^^in her:51.54:51.82^^desperate:51.82:52.16^^attempts to:52.16:52.90^^escape.:52.90:53.94^^The house,:53.94:54.76^^a dark:54.76:55.40^^sanctuary,:55.40:55.92^^swallowed:55.92:56.64^^her back:56.64:57.14^^in, the:57.14:57.84^^walls:57.84:58.02^^whispering:58.02:58.42^^her:58.42:58.70^^secrets,:58.70:59.12^^each more:59.12:59.92^^terrifying:59.92:60.40^^than the:60.40:60.82^^last.:60.82:61.86^^"
		var testText []model.TextDisplay
		textAndTime := strings.Split(textAndTimes, "^^")

		for _, part := range textAndTime {
			textParts := strings.Split(part, ":")
			if len(textParts) < 3 {
				continue
			}
			testText = append(testText, model.TextDisplay{
				Text:      textParts[0],
				StartTime: textParts[1],
				EndTime:   textParts[2],
			})
			audioPath = "./assets/audio/testOutput.wav"
			allText = &testText
		}
	} else {
		// call the api's to get the data dynamically when not in test mode
		videoText, scriptErr := api.GenerateVideoScript(preference.VideoType)
		if scriptErr != nil {
			return model.VideoCreationInfo{}, scriptErr
		}

		audioError := api.GenerateAudioFile(videoText)
		if audioError != nil {
			return model.VideoCreationInfo{}, audioError
		}

		audioPath = "./assets/audio/output.wav"
		foundText, timeStampError := api.TimeStampGenerator(audioPath)
		if timeStampError != nil {
			return model.VideoCreationInfo{}, timeStampError
		}
		allText = foundText
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
	filterFilePath := "filter.txt"

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
	overlayErr := overlayTextOnVideo(preference.FontName, preference.FontColor, combinedPath, outputPath, videoCreationInfo.AllText, filterFilePath)
	if overlayErr != nil {
		return overlayErr
	}
	// delete the audio and video files
	removeErr := deleteFiles(combinedPath, videoCreationInfo.AudioPath, combinedAudioPath, filterFilePath)
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

func overlayTextOnVideo(fontName string, fontColor string, combinedPath string, outputPath string, allText *[]model.TextDisplay, filterFilePath string) error {
	textIntervals := *allText

	// filter the string to place all start and end dates
	fontFile := fmt.Sprintf("./assets/font/%s.ttf", fontName)

	// write the text overlay to a file to prevent having to run a super long command
	filterFile, err := os.Create(filterFilePath)
	if err != nil {
		return fmt.Errorf("error creating the filter.txt file: %v", err)
	}
	defer filterFile.Close()

	for _, textInterval := range textIntervals {
		textFormat := fmt.Sprintf(
			"drawtext=text='%s':fontfile=%s:fontsize=w/6:fontcolor=%s:x=(w-text_w)/2:y=(h-text_h)/2:shadowx=2:shadowy=2:shadowcolor=black:enable='between(t,%s,%s)'",
			textInterval.Text, fontFile, fontColor, textInterval.StartTime, textInterval.EndTime)

		// write the text to the file
		_, err := filterFile.WriteString(textFormat + ",")
		if err != nil {
			return fmt.Errorf("error writing to the filter.txt file: %v", err)
		}
	}
	cmdOverlayText := exec.Command("ffmpeg.exe", "-i", combinedPath, "-filter_complex_script", filterFilePath, "-c:a", "copy", "-y", outputPath)
	overlayOutput, err := cmdOverlayText.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, overlayOutput)
	}

	return nil
}

func deleteFiles(combinedPath string, audioPath string, combinedAudioPath string, filterFilePath string) error {
	combineErr := os.Remove(combinedPath)
	if combineErr != nil {
		return fmt.Errorf("error removing file: %v", combineErr)
	}
	// keep the test audio file
	if audioPath != "./assets/audio/testOutput.wav" {
		outputErr := os.Remove(audioPath)
		if outputErr != nil {
			return fmt.Errorf("error removing file: %v", combineErr)
		}
	}
	combineAudioErr := os.Remove(combinedAudioPath)
	if combineAudioErr != nil {
		return fmt.Errorf("error removing file: %v", combineErr)
	}
	textDisplayErr := os.Remove(filterFilePath)
	if textDisplayErr != nil {
		return fmt.Errorf("error removing file: %v", combineErr)
	}
	return nil
}
