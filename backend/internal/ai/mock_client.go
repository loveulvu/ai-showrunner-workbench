package ai

import (
	"context"
	"fmt"

	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/fidelity"
	"ai-showrunner-workbench/internal/novel"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"
)

type MockClient struct{}

func NewMockClient() MockClient {
	return MockClient{}
}

func (m MockClient) AnalyzeChapter(ctx context.Context, chapter novel.Chapter) (analysis.ChapterAnalysis, error) {
	_ = ctx

	summaryByChapter := map[int]string{
		1: "林澈在雨夜收到失踪父亲留下的铜钥匙，并发现旧剧院线索，人物目标从被动困惑转向主动追查。",
		2: "许岚陪林澈进入旧剧院，二人在地下室发现父亲研究的时钟装置，故事从悬疑线索进入核心设定。",
		3: "顾衡现身索要钥匙，林澈必须决定是否启动时钟装置，主角与阻碍者的冲突正面爆发。",
	}
	summary := summaryByChapter[chapter.Number]
	if summary == "" {
		summary = fmt.Sprintf("第%d章继续推进铜钥匙、旧剧院和父亲失踪之谜。", chapter.Number)
	}

	return analysis.ChapterAnalysis{
		ChapterNumber: chapter.Number,
		ChapterTitle:  chapter.Title,
		Summary:       summary,
		Characters: []analysis.CharacterMention{
			{
				Name:          "林澈",
				RoleInChapter: "追查父亲失踪的主角",
				Traits:        []string{"敏感", "执着", "压抑"},
				StateChange:   "从犹豫接受线索，转为主动靠近旧剧院真相。",
			},
			{
				Name:          "许岚",
				RoleInChapter: "协助林澈调查的盟友",
				Traits:        []string{"谨慎", "理性", "保护欲强"},
				StateChange:   "从担心林澈冒险，转为共同承担调查风险。",
			},
			{
				Name:          "顾衡",
				RoleInChapter: "隐藏实验秘密的阻碍者",
				Traits:        []string{"强势", "克制", "掌控欲强"},
				StateChange:   "从幕后威胁转为正面阻止林澈。",
			},
		},
		Locations: []string{"雨巷", "旧剧院", "地下室", "舞台"},
		KeyEvents: []string{
			fmt.Sprintf("第%d章揭示新的关键线索。", chapter.Number),
			"铜钥匙、剧票和旧时钟被串联为同一条悬疑线。",
		},
		Conflicts: []string{
			"林澈想查清父亲失踪真相，顾衡试图封存旧剧院的记忆实验。",
		},
		FactualAnchors: []string{
			fmt.Sprintf("第%d章标题为%s", chapter.Number, chapter.Title),
			"林澈、许岚、顾衡是示例故事中的核心人物。",
			"铜钥匙、旧剧院和父亲失踪线索是示例故事的关键事实。",
		},
		SceneCandidates: []analysis.SceneCandidate{
			{
				Location:   "雨巷",
				Time:       "夜晚",
				Purpose:    "建立悬疑开端和主角追查动机。",
				Characters: []string{"林澈"},
				KeyEvents:  []string{"林澈收到铜钥匙", "剧票指向旧剧院"},
			},
			{
				Location:   "旧剧院地下室",
				Time:       "深夜",
				Purpose:    "揭示故事核心装置，并让人物关系进入共同冒险。",
				Characters: []string{"林澈", "许岚"},
				KeyEvents:  []string{"二人发现时钟装置", "父亲研究痕迹出现"},
			},
			{
				Location:   "旧剧院舞台",
				Time:       "黎明前",
				Purpose:    "让主角与阻碍者正面对峙，推动第一幕高潮。",
				Characters: []string{"林澈", "许岚", "顾衡"},
				KeyEvents:  []string{"顾衡索要钥匙", "林澈决定启动时钟"},
			},
		},
	}, nil
}

func (m MockClient) MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (story.StoryBible, error) {
	_ = ctx

	timeline := make([]story.TimelineEvent, 0, len(analyses))
	for _, item := range analyses {
		event := item.Summary
		if len(item.KeyEvents) > 0 {
			event = item.KeyEvents[0]
		}
		timeline = append(timeline, story.TimelineEvent{
			ChapterNumber: item.ChapterNumber,
			Event:         event,
		})
	}

	return story.StoryBible{
		Title:   "雨夜旧剧院",
		Logline: "一个青年追查父亲失踪之谜，却在旧剧院发现能改写记忆的时钟装置。",
		GlobalCharacters: []story.Character{
			{ID: "char_lin_che", Name: "林澈", Role: "主角", Motivation: "查清父亲失踪真相"},
			{ID: "char_xu_lan", Name: "许岚", Role: "盟友", Motivation: "保护林澈并揭开剧院秘密"},
			{ID: "char_gu_heng", Name: "顾衡", Role: "阻碍者", Motivation: "夺回铜钥匙并隐藏旧实验"},
		},
		Timeline:     timeline,
		MainConflict: "林澈要公开父亲失踪真相，顾衡要封存剧院里的记忆实验。",
		ScenePlan: []story.ScenePlanItem{
			{ID: "scene_001", SourceChapter: 1, Summary: "林澈在雨巷收到父亲留下的铜钥匙。", Location: "雨巷", Time: "夜晚", Characters: []string{"char_lin_che"}},
			{ID: "scene_002", SourceChapter: 2, Summary: "林澈和许岚进入旧剧院地下室，发现时钟装置。", Location: "旧剧院地下室", Time: "深夜", Characters: []string{"char_lin_che", "char_xu_lan"}},
			{ID: "scene_003", SourceChapter: 3, Summary: "顾衡逼迫林澈交出钥匙，三人在舞台上对峙。", Location: "旧剧院舞台", Time: "黎明前", Characters: []string{"char_lin_che", "char_xu_lan", "char_gu_heng"}},
		},
	}, nil
}

func (m MockClient) GenerateScreenplay(ctx context.Context, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (screenplay.Screenplay, error) {
	_ = ctx

	chapterTitles := map[int]string{
		1: "雨夜钥匙",
		2: "旧剧院",
		3: "舞台对峙",
	}
	sourceChapters := make([]screenplay.SourceChapter, 0, len(analyses))
	for _, item := range analyses {
		sourceChapters = append(sourceChapters, screenplay.SourceChapter{
			Number:  item.ChapterNumber,
			Title:   item.ChapterTitle,
			Summary: item.Summary,
		})
	}
	if len(sourceChapters) == 0 {
		sourceChapters = make([]screenplay.SourceChapter, 0, len(bible.Timeline))
	}
	for _, item := range bible.Timeline {
		if len(analyses) > 0 {
			break
		}
		title := chapterTitles[item.ChapterNumber]
		if title == "" {
			title = fmt.Sprintf("第%d章", item.ChapterNumber)
		}
		sourceChapters = append(sourceChapters, screenplay.SourceChapter{
			Number:  item.ChapterNumber,
			Title:   title,
			Summary: item.Event,
		})
	}

	characterDescriptions := map[string]string{
		"char_lin_che": "背负父亲失踪阴影的青年，外表克制，内心渴望确认父亲留下的真相。",
		"char_xu_lan":  "林澈的朋友和调查伙伴，负责把危险拉回现实，也推动林澈说出真实恐惧。",
		"char_gu_heng": "旧剧院实验的守门人，试图阻止记忆装置再次启动。",
	}
	characters := make([]screenplay.Character, 0, len(bible.GlobalCharacters))
	for _, character := range bible.GlobalCharacters {
		characters = append(characters, screenplay.Character{
			ID:          character.ID,
			Name:        character.Name,
			Role:        character.Role,
			Description: characterDescriptions[character.ID],
		})
	}

	scenes := make([]screenplay.Scene, 0, len(bible.ScenePlan))
	for _, plan := range bible.ScenePlan {
		scenes = append(scenes, screenplay.Scene{
			ID:            plan.ID,
			SourceChapter: plan.SourceChapter,
			Location:      plan.Location,
			Time:          plan.Time,
			Summary:       plan.Summary,
			Characters:    plan.Characters,
			Dialogues: []screenplay.Dialogue{
				{Character: "char_lin_che", Emotion: "克制", Line: "这把钥匙不是遗物，它像是在等我回来。"},
				{Character: "char_xu_lan", Emotion: "紧张", Line: "如果这里真的能保存记忆，我们就不能让顾衡先找到它。"},
			},
			Actions: []string{
				"雨水顺着破损的屋檐滴落，铜钥匙在林澈掌心微微发亮。",
				"远处的舞台灯忽明忽暗，像有人刚刚经过。",
			},
		})
	}

	return screenplay.Screenplay{
		Title:          bible.Title,
		SourceChapters: sourceChapters,
		Characters:     characters,
		Scenes:         scenes,
	}, nil
}

func (m MockClient) CheckFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (fidelity.FidelityResult, error) {
	_ = ctx
	_ = current
	_ = bible
	_ = analyses
	return fidelity.FidelityResult{Passed: true, Issues: []fidelity.FidelityIssue{}}, nil
}

func (m MockClient) RepairFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis, result fidelity.FidelityResult) (screenplay.Screenplay, error) {
	_ = ctx
	_ = bible
	_ = analyses
	_ = result
	return current, nil
}
