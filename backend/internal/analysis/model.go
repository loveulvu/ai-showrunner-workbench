package analysis

import "encoding/json"

type CharacterMention struct {
	Name          string   `json:"name"`
	RoleInChapter string   `json:"role_in_chapter"`
	Traits        []string `json:"traits"`
	StateChange   string   `json:"state_change"`
}

func (c *CharacterMention) UnmarshalJSON(data []byte) error {
	type rawCharacterMention struct {
		Name          string          `json:"name"`
		RoleInChapter string          `json:"role_in_chapter"`
		Traits        json.RawMessage `json:"traits"`
		StateChange   string          `json:"state_change"`
	}

	var raw rawCharacterMention
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.Name = raw.Name
	c.RoleInChapter = raw.RoleInChapter
	c.Traits = unmarshalStringList(raw.Traits)
	c.StateChange = raw.StateChange
	return nil
}

type SceneCandidate struct {
	Location   string   `json:"location"`
	Time       string   `json:"time"`
	Purpose    string   `json:"purpose"`
	Characters []string `json:"characters"`
	KeyEvents  []string `json:"key_events"`
}

type ChapterAnalysis struct {
	ChapterNumber   int                `json:"chapter_number"`
	ChapterTitle    string             `json:"chapter_title"`
	Summary         string             `json:"summary"`
	Characters      []CharacterMention `json:"characters"`
	Locations       []string           `json:"locations"`
	KeyEvents       []string           `json:"key_events"`
	Conflicts       []string           `json:"conflicts"`
	SceneCandidates []SceneCandidate   `json:"scene_candidates"`
	FactualAnchors  []string           `json:"factual_anchors"`
}

func unmarshalStringList(data json.RawMessage) []string {
	if len(data) == 0 || string(data) == "null" {
		return []string{}
	}

	var values []string
	if err := json.Unmarshal(data, &values); err == nil {
		if values == nil {
			return []string{}
		}
		return values
	}

	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		if value == "" {
			return []string{}
		}
		return []string{value}
	}

	return []string{}
}
