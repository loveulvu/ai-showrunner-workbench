package editor

type ClipAsset struct {
	ShotID          string `json:"shot_id"`
	SourceURL       string `json:"source_url"`
	LocalPath       string `json:"local_path"`
	DurationSeconds int    `json:"duration_seconds"`
	Subtitle        string `json:"subtitle"`
}

type EditingPlan struct {
	Clips       []ClipAsset `json:"clips"`
	OutputFile  string      `json:"output_file"`
	AspectRatio string      `json:"aspect_ratio"`
	Resolution  string      `json:"resolution"`
	FPS         int         `json:"fps"`
}

type EditResult struct {
	OutputFile    string   `json:"output_file"`
	ClipCount     int      `json:"clip_count"`
	SubtitlesFile string   `json:"subtitles_file"`
	Warnings      []string `json:"warnings"`
}
