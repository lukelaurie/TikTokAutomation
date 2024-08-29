package database

import (
	// "database/sql"
	"database/sql"
	"fmt"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func RetrieveSchedulerInfo(username string, curScheduler int) (model.Preference, error) {
	var preference model.Preference

	err := DB.First(&preference, "username = ? AND preference_order = ?", username, curScheduler).Error
	if err != nil {
		// check if the error was from no user being found
		if err == sql.ErrNoRows {
			return preference, fmt.Errorf("username in preference or preference index is invalid")
		}
		return preference, err
	}

	return preference, nil
}

func RetrievePreferenceTracker(username string) (model.PreferenceTracker, error) {
	var preferenceTracker model.PreferenceTracker

	err := DB.First(&preferenceTracker, "username = ?", username).Error
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
	newPreference := model.Preference{
		Username: username,
		VideoType: preference.VideoType,
		BackgroundVideoType: preference.BackgroundVideoType,
		FontName: preference.FontName,
		FontColor: preference.FontColor,
		PreferenceOrder: preferenceIndex,
	}

	// insert the user into the database 
	err := DB.Create(&newPreference).Error
	return err
}

func AddNewUserPreferenceTracker(username string) error {
	preferenceTracker := model.PreferenceTracker{
		Username: username,
	}

	// insert the user into the database 
	err := DB.Create(&preferenceTracker).Error
	return err
}

func IncrementPreferenceTracker(username string, preferenceIndex int, isPreferenceCount bool) error {
	var preferenceTracker model.PreferenceTracker
	err := DB.First(&preferenceTracker, "username = ?", username).Error
	if err != nil {
		return err
	} 

	// determine if increasing total count of preferences or changing the index pointing to
	if isPreferenceCount {
		preferenceTracker.CurrentPreferenceCount += 1
	} else {
		preferenceTracker.CurrentPreferenceOrder = preferenceIndex
	}

	// update the preference tracker with the new count
	err = DB.Save(&preferenceTracker).Error
	return err
}
