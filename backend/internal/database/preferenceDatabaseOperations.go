package database

import (
	"database/sql"
	"fmt"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func RetrieveSchedulerInfo(username string, curScheduler int) (model.Preference, error) {
	var preference model.Preference

	query := `SELECT video_type, background_video_type, font_name, font_color FROM preferences WHERE username = $1 AND preference_order = $2`

	// search for the row
	err := DB.QueryRow(query, username, curScheduler).Scan(&preference.VideoType, &preference.BackgroundVideoType,
		&preference.FontName, &preference.FontColor)
	if err != nil {
		// check if the error was from no user being found
		if err == sql.ErrNoRows {
			return preference, fmt.Errorf("username in preference or preference index is invalid")
		}
		return preference, err
	}
	return preference, nil
}

func RetrieveAndUpdatePreferenceIndex(username string) (int, error) {
	// get the current index for the preference
	preferenceTracker, err := RetrievePreferenceTracker(username)
	if err != nil {
		return -1, err
	}

	if preferenceTracker.CurPreferenceCount == 0 {
		return -1, fmt.Errorf("username has not yet created a preference")
	}

	// increment the preference to next index or restart to front if already created a video for all
	newPreferenceIndex := preferenceTracker.CurPreferenceOrder + 1
	if newPreferenceIndex > preferenceTracker.CurPreferenceCount {
		newPreferenceIndex = 1
	}

	return preferenceTracker.CurPreferenceOrder, nil
}
func RetrievePreferenceTracker(username string) (model.PreferenceTracker, error) {
	var preferenceTracker model.PreferenceTracker

	query := `SELECT id, username, current_preference_order, current_preference_count FROM user_preference_tracker WHERE username = $1`

	// search for the row
	err := DB.QueryRow(query, username).Scan(&preferenceTracker.Id, &preferenceTracker.Username,
		&preferenceTracker.CurPreferenceOrder, &preferenceTracker.CurPreferenceCount)
	if err != nil {
		// check if the error was from no user being found
		if err == sql.ErrNoRows {
			return preferenceTracker, fmt.Errorf("username in preference tacker table invalid")
		}
		return preferenceTracker, err
	}
	return preferenceTracker, nil
}

func AddNewUserPreference(preference model.Preference, username string, preferenceIndex int) error {
	query := `INSERT INTO preferences (username, video_type, background_video_type, 
		font_name, font_color, preference_order) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := DB.Exec(query, username, preference.VideoType, preference.BackgroundVideoType,
		preference.FontName, preference.FontColor, preferenceIndex)
	return err
}

func AddNewUserPreferenceTracker(username string) error {
	query := `INSERT INTO user_preference_tracker (username) VALUES ($1)`
	_, err := DB.Exec(query, username)
	return err
}

func IncrementPreferenceTracker(username string, preferenceIndex int, isPreferenceCount bool) error {
	var preferenceCount string
	if isPreferenceCount {
		preferenceCount = "current_preference_count"
	} else {
		preferenceCount = "current_preference_order"
	}

	query := fmt.Sprintf(`UPDATE user_preference_tracker SET %s=%d WHERE username=($1)`, preferenceCount, preferenceIndex)
	_, err := DB.Exec(query, username)
	return err
}
