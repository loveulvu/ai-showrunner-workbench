package handlers

import (
	"log"
	"net/http"
	"unicode/utf8"

	"ai-showrunner-workbench/internal/ai"
	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/fidelity"
	"ai-showrunner-workbench/internal/novel"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"

	"github.com/gin-gonic/gin"
)

type generateRequest struct {
	NovelText string `json:"novel_text"`
}

type generateResponse struct {
	ChapterCount     int                         `json:"chapter_count"`
	ChapterAnalyses  []analysis.ChapterAnalysis  `json:"chapter_analyses"`
	StoryBible       story.StoryBible            `json:"story_bible"`
	ScreenplayJSON   screenplay.Screenplay       `json:"screenplay_json"`
	ScreenplayYAML   string                      `json:"screenplay_yaml"`
	ValidationResult screenplay.ValidationResult `json:"validation"`
	FidelityResult   fidelity.FidelityResult     `json:"fidelity_result"`
	Meta             generateMeta                `json:"meta"`
}

type generateMeta struct {
	AIProvider string `json:"ai_provider"`
	AIModel    string `json:"ai_model"`
}

type parserDebug struct {
	NovelTextLength     int      `json:"novel_text_length"`
	ParsedChapterCount  int      `json:"parsed_chapter_count"`
	ParsedChapterTitles []string `json:"parsed_chapter_titles"`
	First300Chars       string   `json:"first_300_chars"`
}

func Generate(c *gin.Context) {

	var req generateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json body",
		})
		return
	}

	chapters := novel.ParseChapters(req.NovelText)
	titles := parsedChapterTitles(chapters)
	log.Printf("parser debug: novel_text_length=%d parsed_chapter_count=%d parsed_chapter_titles=%q", utf8.RuneCountInString(req.NovelText), len(chapters), titles)
	if len(chapters) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "at least 3 chapters are required",
			"debug": parserDebug{
				NovelTextLength:     utf8.RuneCountInString(req.NovelText),
				ParsedChapterCount:  len(chapters),
				ParsedChapterTitles: titles,
				First300Chars:       firstNChars(req.NovelText, 300),
			},
		})
		return
	}
	client, err := ai.NewClientFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	runtimeStatus := ai.RuntimeStatusFromEnv()
	analyzer := analysis.NewAnalyzer(client)
	merger := story.NewMerger(client)
	generator := screenplay.NewGenerator(client)
	fidelityChecker := fidelity.NewChecker(client)

	ctx := c.Request.Context()
	chapterAnalyses, err := analyzer.AnalyzeChapters(ctx, chapters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	storyBible, err := merger.Merge(ctx, chapterAnalyses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	screenplayJSON, err := generator.Generate(ctx, storyBible, chapterAnalyses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	applySourceChaptersFromAnalyses(&screenplayJSON, chapterAnalyses)

	screenplayJSON, fidelityResult, err := fidelityChecker.CheckAndRepairOnce(ctx, screenplayJSON, storyBible, chapterAnalyses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	applySourceChaptersFromAnalyses(&screenplayJSON, chapterAnalyses)

	validation := screenplay.Validate(screenplayJSON)
	screenplayYAML := screenplay.ToYAML(screenplayJSON)

	c.JSON(http.StatusOK, generateResponse{
		ChapterCount:     len(chapters),
		ChapterAnalyses:  chapterAnalyses,
		StoryBible:       storyBible,
		ScreenplayJSON:   screenplayJSON,
		ScreenplayYAML:   screenplayYAML,
		ValidationResult: validation,
		FidelityResult:   fidelityResult,
		Meta: generateMeta{
			AIProvider: runtimeStatus.AIProvider,
			AIModel:    runtimeStatus.AIModel,
		},
	})
}

func applySourceChaptersFromAnalyses(target *screenplay.Screenplay, analyses []analysis.ChapterAnalysis) {
	sourceChapters := make([]screenplay.SourceChapter, 0, len(analyses))
	for _, chapter := range analyses {
		sourceChapters = append(sourceChapters, screenplay.SourceChapter{
			Number:  chapter.ChapterNumber,
			Title:   chapter.ChapterTitle,
			Summary: chapter.Summary,
		})
	}
	target.SourceChapters = sourceChapters
}

func parsedChapterTitles(chapters []novel.Chapter) []string {
	titles := make([]string, 0, len(chapters))
	for _, chapter := range chapters {
		titles = append(titles, chapter.Title)
	}
	return titles
}

func firstNChars(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}
