package storage

type Preferences struct {
	ThemeID         string `json:"theme_id"`
	Mode            string `json:"mode"`
	DurationSeconds int    `json:"duration_seconds"`
	WordCount       int    `json:"word_count"`
	Punctuation     bool   `json:"punctuation"`
	Numbers         bool   `json:"numbers"`
}

type BestScore struct {
	WPM       int   `json:"wpm"`
	Accuracy  int   `json:"accuracy"`
	Timestamp int64 `json:"timestamp"`
}

type Data struct {
	Preferences Preferences          `json:"preferences"`
	BestScores  map[string]BestScore `json:"best_scores"`
}
