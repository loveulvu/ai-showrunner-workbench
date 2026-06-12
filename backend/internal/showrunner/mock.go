package showrunner

const defaultVideoStyle = "cinematic xianxia short drama"
const defaultNegativePrompt = "blurry, distorted face, inconsistent character, different outfit, extra limbs, bad hands, text, subtitles, watermark, logo, low quality, jump cut, flickering, deformed body"

func MockResult(input GenerateInput) ShowrunnerResult {
	input = PrepareInput(input)
	title := input.Screenplay.Title
	if title == "" {
		title = input.StoryBible.Title
	}
	style := input.Style
	if style == "" {
		style = defaultVideoStyle
	}

	leadIdentity := CharacterVisualIdentity{
		Age: "24", Face: "oval face, defined brows, focused dark eyes", Hairstyle: "long black hair in a high half-ponytail",
		Costume: "midnight-blue fitted xianxia robe with silver trim and dark leather belt", ColorPalette: "midnight blue, silver, charcoal",
		ExpressionBaseline: "controlled alertness", BodyType: "lean athletic", ConsistencyPrompt: "same oval face, high half-ponytail, midnight-blue silver-trim robe, lean build",
	}
	partnerIdentity := CharacterVisualIdentity{
		Age: "26", Face: "angular face, calm dark eyes", Hairstyle: "black hair tied in a low practical knot",
		Costume: "charcoal layered xianxia robe with muted teal sash", ColorPalette: "charcoal, muted teal, black",
		ExpressionBaseline: "watchful concern", BodyType: "tall balanced build", ConsistencyPrompt: "same angular face, low hair knot, charcoal robe and muted teal sash, tall build",
	}
	openingIdentity := SceneVisualIdentity{
		Architecture: "weathered timber corridor beside an old theatre courtyard", Lighting: "cool moonlight from frame left with warm lantern rim light",
		ColorPalette: "deep blue, charcoal, restrained amber", Atmosphere: "tense quiet night after rain", KeyProps: FlexibleStringList{"bronze key", "old theatre ticket", "hanging lantern"},
		ConsistencyPrompt: "same wet timber corridor, lantern positions, moonlight direction, bronze key and old ticket",
	}

	return ShowrunnerResult{
		Characters: []CharacterProfile{
			{ID: "char_lead", Name: "Lead", Role: "protagonist", Personality: FlexibleStringList{"determined", "observant"}, Appearance: FlexibleStringList{leadIdentity.ConsistencyPrompt}, Costume: FlexibleStringList{leadIdentity.Costume}, VoiceStyle: FlexibleStringList{"calm with restrained urgency"}, KeyMotivation: FlexibleStringList{"discover the truth"}, ConsistencyNotes: FlexibleStringList{"keep face, hair, costume, and proportions unchanged"}, VisualIdentity: leadIdentity},
			{ID: "char_partner", Name: "Partner", Role: "ally", Personality: FlexibleStringList{"careful", "supportive"}, Appearance: FlexibleStringList{partnerIdentity.ConsistencyPrompt}, Costume: FlexibleStringList{partnerIdentity.Costume}, VoiceStyle: FlexibleStringList{"grounded and direct"}, KeyMotivation: FlexibleStringList{"protect the lead"}, ConsistencyNotes: FlexibleStringList{"keep face, hair, costume, and proportions unchanged"}, VisualIdentity: partnerIdentity},
		},
		Scenes: []SceneProfile{
			{ID: "scene_001", Name: "Old Theatre Corridor", Location: "old theatre courtyard corridor", TimeOfDay: "night", Atmosphere: openingIdentity.Atmosphere, VisualStyle: style, KeyProps: openingIdentity.KeyProps, ConsistencyNotes: FlexibleStringList{"preserve layout, light direction, rain level, and prop placement"}, VisualIdentity: openingIdentity},
			{ID: "scene_002", Name: "Old Theatre Stage", Location: "abandoned theatre stage", TimeOfDay: "night", Atmosphere: "ominous stillness", VisualStyle: style, KeyProps: FlexibleStringList{"dusty curtain", "wooden stage clock"}, ConsistencyNotes: FlexibleStringList{"preserve stage geometry and clock position"}, VisualIdentity: SceneVisualIdentity{Architecture: "aged timber stage with deep red curtain", Lighting: "single cool overhead shaft with dim amber footlights", ColorPalette: "charcoal, faded red, restrained amber", Atmosphere: "ominous stillness", KeyProps: FlexibleStringList{"dusty curtain", "wooden stage clock"}, ConsistencyPrompt: "same aged timber stage, faded red curtain, overhead light shaft, and wooden stage clock"}},
		},
		Chapters: []ChapterBreakdown{
			{ChapterNumber: 1, ChapterTitle: title, Summary: "The lead and partner discover that the bronze key reacts to the old theatre ticket.", MainCharacters: FlexibleStringList{"char_lead", "char_partner"}, MainScenes: FlexibleStringList{"scene_001"}, EmotionalArc: "caution to revelation", KeyEvents: FlexibleStringList{"enter corridor", "key activates", "partners react"}},
		},
		Shots: []Shot{
			{
				ID: "shot_001", ChapterNumber: 1, SceneID: "scene_001", Characters: FlexibleStringList{"char_lead", "char_partner"},
				Action: "The lead enters the wet corridor holding the bronze key while the partner follows and stops one step behind.", Camera: "medium-wide tracking shot", Background: openingIdentity.ConsistencyPrompt, DurationHint: "5s",
				CharacterVisuals: FlexibleStringList{leadIdentity.ConsistencyPrompt, partnerIdentity.ConsistencyPrompt}, SceneVisuals: openingIdentity.ConsistencyPrompt,
				CameraAngle: "eye-level medium-wide", CameraMovement: "slow controlled lateral track", Composition: "lead foreground left, partner midground right, lantern depth lines", Lighting: openingIdentity.Lighting,
				Motion: "measured footsteps, robe hems and damp hair move subtly, lantern flame reacts to their passage", ContinuityNotes: "establish positions, costumes, key in lead's right hand, and left-to-right screen direction for shots 2 and 3",
				ImagePrompt:    style + ", same lead and partner identities, wet old theatre corridor, bronze key, cinematic medium-wide frame",
				VideoPrompt:    style + ", same lead oval face and high half-ponytail in midnight-blue silver-trim robe, same partner angular face and low hair knot in charcoal robe with teal sash, wet old theatre corridor, lead enters holding bronze key as partner follows, slow lateral tracking, cool moonlight and warm lantern rim, medium-wide layered composition, stable faces, stable costumes, realistic walking motion, no text, no subtitles, no watermark",
				NegativePrompt: defaultNegativePrompt, AudioPrompt: "soft footsteps on wet timber, low night ambience",
			},
			{
				ID: "shot_002", ChapterNumber: 1, SceneID: "scene_001", Characters: FlexibleStringList{"char_lead", "char_partner"},
				Action: "Without changing position, the lead raises the bronze key over the old ticket; a restrained blue glow appears and the partner leans closer.", Camera: "medium close two-shot", Background: openingIdentity.ConsistencyPrompt, DurationHint: "5s",
				CharacterVisuals: FlexibleStringList{leadIdentity.ConsistencyPrompt, partnerIdentity.ConsistencyPrompt}, SceneVisuals: openingIdentity.ConsistencyPrompt,
				CameraAngle: "slightly low medium close two-shot", CameraMovement: "gentle push-in toward key and faces", Composition: "glowing key centered between both faces, lead left and partner right", Lighting: openingIdentity.Lighting + ", subtle blue key glow",
				Motion: "lead slowly raises key, partner leans in, eyes track the glow, sleeves settle naturally", ContinuityNotes: "continue exact positions, costumes, key hand, screen direction, corridor, and lighting from shot 1",
				ImagePrompt:    style + ", same characters and costumes, bronze key glowing above old ticket, cinematic close two-shot",
				VideoPrompt:    style + ", same lead face hair and midnight-blue silver-trim robe, same partner face hair and charcoal teal-sash robe, same wet theatre corridor and lanterns, lead raises bronze key above old ticket as blue glow appears and partner leans closer, gentle camera push-in, cool moonlight warm rim and subtle blue glow, centered two-shot composition, stable faces, stable costumes, realistic hand and eye motion, no text, no subtitles, no watermark",
				NegativePrompt: defaultNegativePrompt, AudioPrompt: "quiet magical hum under restrained night ambience",
			},
			{
				ID: "shot_003", ChapterNumber: 1, SceneID: "scene_001", Characters: FlexibleStringList{"char_lead", "char_partner"}, Dialogue: FlexibleStringList{"It is pointing inside."},
				Action: "The glow fades; the lead looks toward the theatre door while the partner reacts with concern and grips the lead's sleeve.", Camera: "close reaction two-shot", Background: openingIdentity.ConsistencyPrompt, DurationHint: "5s",
				CharacterVisuals: FlexibleStringList{leadIdentity.ConsistencyPrompt, partnerIdentity.ConsistencyPrompt}, SceneVisuals: openingIdentity.ConsistencyPrompt,
				CameraAngle: "eye-level close reaction shot", CameraMovement: "controlled rack focus from fading key to faces", Composition: "key low foreground, lead profile left, partner reaction right, theatre door in background", Lighting: openingIdentity.Lighting + ", fading blue reflection",
				Motion: "key glow fades, lead turns eyes then head toward door, partner grips sleeve and exhales", ContinuityNotes: "preserve positions, wardrobe, key hand, wet corridor, lanterns, and screen direction from shots 1 and 2",
				ImagePrompt:    style + ", same characters and costumes reacting after key glow fades, theatre door behind them",
				VideoPrompt:    style + ", same lead face high half-ponytail and midnight-blue silver-trim robe, same partner angular face low knot and charcoal teal-sash robe, same wet theatre corridor and props, glow fades as lead turns toward theatre door and partner grips the lead's sleeve, controlled rack focus from key to faces, cool moonlight warm lantern rim, close reaction composition, stable faces, stable costumes, realistic reaction motion, no text, no subtitles, no watermark",
				NegativePrompt: defaultNegativePrompt, AudioPrompt: "glow fades into silence, fabric grip, distant wooden creak",
			},
		},
		AssetPrompts: AssetPromptSet{
			CharacterPrompts:  map[string]string{"char_lead": leadIdentity.ConsistencyPrompt, "char_partner": partnerIdentity.ConsistencyPrompt},
			BackgroundPrompts: map[string]string{"scene_001": openingIdentity.ConsistencyPrompt},
			ShotPrompts:       map[string]string{"shot_001": "continuous setup entrance", "shot_002": "continuous key activation", "shot_003": "continuous reaction and result"},
			VoicePrompts:      map[string]string{"char_lead": "calm voice with restrained urgency", "char_partner": "grounded supportive voice"},
		},
		Warnings: FlexibleStringList{},
		Mode:     input.Mode,
	}
}
