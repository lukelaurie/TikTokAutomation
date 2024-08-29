package model

type User struct {
	Username string `gorm:"primaryKey"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

type Preference struct {
	ID                 int    `gorm:"primaryKey;autoIncrement"`
	Username           string `gorm:"not null;index;references:Username"`
	VideoType          string `gorm:"not null"`
	BackgroundVideoType string `gorm:"not null"`
	FontName           string `gorm:"not null"`
	FontColor          string `gorm:"not null"`
	PreferenceOrder    int    
}

type PreferenceTracker struct {
	ID                    int    `gorm:"primaryKey;autoIncrement"`
	Username              string `gorm:"not null;index;references:Username"`
	CurrentPreferenceOrder int    `gorm:"default:1"`
	CurrentPreferenceCount int    `gorm:"default:0"`
}