package novel

type Chapter struct {
	Number int    `json:"chapter_number"`
	Title  string `json:"chapter_title"`
	Text   string `json:"text"`
}
