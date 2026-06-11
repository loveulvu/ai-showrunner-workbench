package story

type Character struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Motivation string `json:"motivation"`
}

type TimelineEvent struct {
	ChapterNumber int    `json:"chapter_number"`
	Event         string `json:"event"`
}

type ScenePlanItem struct {
	ID            string   `json:"id"`
	SourceChapter int      `json:"source_chapter"`
	Summary       string   `json:"summary"`
	Location      string   `json:"location"`
	Time          string   `json:"time"`
	Characters    []string `json:"characters"`
}

type StoryBible struct {
	Title            string          `json:"title"`
	Logline          string          `json:"logline"`
	GlobalCharacters []Character     `json:"global_characters"`
	Timeline         []TimelineEvent `json:"timeline"`
	MainConflict     string          `json:"main_conflict"`
	ScenePlan        []ScenePlanItem `json:"scene_plan"`
}
