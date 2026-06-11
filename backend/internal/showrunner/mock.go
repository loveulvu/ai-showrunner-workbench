package showrunner

func MockResult(input GenerateInput) ShowrunnerResult {
	title := input.Screenplay.Title
	if title == "" {
		title = input.StoryBible.Title
	}

	return ShowrunnerResult{
		Characters: []CharacterProfile{
			{ID: "char_lead", Name: "Lead", Role: "protagonist", Personality: []string{"determined", "observant"}, Appearance: "recognizable silhouette and focused expression", Costume: "consistent practical hero costume", VoiceStyle: "calm with restrained urgency", KeyMotivation: "discover the truth", ConsistencyNotes: []string{"keep costume colors and facial features consistent"}},
			{ID: "char_partner", Name: "Partner", Role: "ally", Personality: []string{"careful", "supportive"}, Appearance: "clear supporting-character silhouette", Costume: "consistent layered supporting costume", VoiceStyle: "grounded and direct", KeyMotivation: "protect the lead", ConsistencyNotes: []string{"maintain the same costume and proportions"}},
		},
		Scenes: []SceneProfile{
			{ID: "scene_001", Name: "Opening Location", Location: "story opening location", TimeOfDay: "night", Atmosphere: "tense and mysterious", VisualStyle: input.Style, KeyProps: []string{"story clue"}, ConsistencyNotes: []string{"reuse the same spatial layout"}},
			{ID: "scene_002", Name: "Confrontation Location", Location: "main confrontation location", TimeOfDay: "late night", Atmosphere: "dramatic and focused", VisualStyle: input.Style, KeyProps: []string{"central story object"}, ConsistencyNotes: []string{"preserve lighting direction and prop placement"}},
		},
		Chapters: []ChapterBreakdown{
			{ChapterNumber: 1, ChapterTitle: title, Summary: "The lead discovers a clue and moves toward the central conflict.", MainCharacters: []string{"char_lead", "char_partner"}, MainScenes: []string{"scene_001", "scene_002"}, EmotionalArc: "uncertainty to resolve", KeyEvents: []string{"clue discovered", "conflict approached"}},
		},
		Shots: []Shot{
			{ID: "shot_001", ChapterNumber: 1, SceneID: "scene_001", Characters: []string{"char_lead"}, Action: "The lead studies the clue in the opening location.", Camera: "medium close-up, slow push in", Background: "consistent opening location", DurationHint: "4s", ImagePrompt: "cinematic animation frame, lead studying a clue, consistent character design", VideoPrompt: "slow push-in while the lead studies the clue", AudioPrompt: "quiet ambience with restrained tension"},
			{ID: "shot_002", ChapterNumber: 1, SceneID: "scene_001", Characters: []string{"char_lead", "char_partner"}, Dialogue: "We need to understand what this means.", Camera: "two-shot at eye level", Background: "consistent opening location", DurationHint: "5s", ImagePrompt: "cinematic animation two-shot, lead and partner discussing a clue", VideoPrompt: "subtle character movement during dialogue", AudioPrompt: "clear dialogue with low ambient room tone"},
			{ID: "shot_003", ChapterNumber: 1, SceneID: "scene_002", Characters: []string{"char_lead", "char_partner"}, Action: "They enter the confrontation location and face the central conflict.", Camera: "wide establishing shot", Background: "consistent confrontation location", DurationHint: "6s", ImagePrompt: "wide cinematic animation frame, characters entering confrontation location", VideoPrompt: "", AudioPrompt: ""},
		},
		AssetPrompts: AssetPromptSet{
			CharacterPrompts:  map[string]string{"char_lead": "consistent animated protagonist character sheet", "char_partner": "consistent animated ally character sheet"},
			BackgroundPrompts: map[string]string{"scene_001": "reusable opening-location background", "scene_002": "reusable confrontation-location background"},
			ShotPrompts:       map[string]string{"shot_001": "lead studies clue", "shot_002": "two-character dialogue", "shot_003": "wide confrontation entrance"},
			VoicePrompts:      map[string]string{"char_lead": "calm voice with restrained urgency", "char_partner": "grounded supportive voice"},
		},
		Warnings: []string{},
	}
}
