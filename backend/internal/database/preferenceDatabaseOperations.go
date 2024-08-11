package database

import (
	"database/sql"
	"fmt"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func RetrieveSchedulerInfo(curScheduler int) (string, string, string, string) {
	// pull out the list of schedues from the member struct

	// get the specific schedule needed
	return "scary", "minecraft", "Proxima-Nova-Semibold", "#FFFFFF"
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
		font_name, font_color, preference_order) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := DB.Exec(query, username, preference, preference.VideoType, preference.BackgroundVideoType,
		preference.FontName, preference.FontColor, preferenceIndex)
	return err
}

func AddNewUserPreferenceTracker(username string) error {
	query := `INSERT INTO user_preference_tracker (username) VALUES ($1)`
	_, err := DB.Exec(query, username)
	return err
}

func IncrementPreferenceCount (preference model.Preference, username string, preferenceIndex int) error {
	query := `INSERT INTO user_preference_tracker (username) VALUES ($1)`
	_, err := DB.Exec(query, username)
	return err
}
