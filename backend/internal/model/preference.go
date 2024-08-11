package model

type Preference struct {
	VideoType           string `json:"videoType"`
	BackgroundVideoType string `json:"backgroundVideoType"`
	FontName            string `json:"fontName"`
	FontColor           string `json:"fontColor"`
}

type PreferenceTracker struct {
	Id                 string `json:"id"`
	Username           string `json:"username"`
	CurPreferenceOrder int `json:"current_preference_order"`
	CurPreferenceCount int `json:"current_preference_count"`
}
