package novel

import (
	"strings"
	"testing"
)

func TestParseHeadingsOnSeparateLines(t *testing.T) {
	input := `第四节：古月方源！
正文A
第五节：人祖三蛊，希望开窍
正文B
第六节：未来的路，会很精彩
正文C`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "古月方源！", "正文A")
	assertChapter(t, chapters[1], 2, "人祖三蛊，希望开窍", "正文B")
	assertChapter(t, chapters[2], 3, "未来的路，会很精彩", "正文C")
}

func TestParseChineseChapters(t *testing.T) {
	input := `第1章 开端
这里是开端正文
第二章 发展
这里是发展正文
第十二章 雨夜
这里是雨夜正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
	assertChapter(t, chapters[2], 3, "雨夜", "这里是雨夜正文")
}

func TestParseChaptersChineseHeadingsWithoutSpaces(t *testing.T) {
	input := `第1章开端
这里是开端正文
第一章发展
这里是发展正文
第十二章雨夜
这里是雨夜正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
	assertChapter(t, chapters[2], 3, "雨夜", "这里是雨夜正文")
}

func TestParseChaptersChineseSectionHeadings(t *testing.T) {
	input := `第一节：相遇
这里是相遇正文
第二节 追踪
这里是追踪正文
第三节对峙
这里是对峙正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "相遇", "这里是相遇正文")
	assertChapter(t, chapters[1], 2, "追踪", "这里是追踪正文")
	assertChapter(t, chapters[2], 3, "对峙", "这里是对峙正文")
}

func TestParseEnglishChapters(t *testing.T) {
	input := `Chapter 1 Beginning
Chapter one body
chapter 2 Middle
Chapter two body
Chapter 3: Ending
Chapter three body`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "Beginning", "Chapter one body")
	assertChapter(t, chapters[1], 2, "Middle", "Chapter two body")
	assertChapter(t, chapters[2], 3, "Ending", "Chapter three body")
}

func TestParseChaptersIgnoresPreface(t *testing.T) {
	input := `前言
这是一段简介，不应该变成未命名章节。

第1章 开端
这里是开端正文
第二章 发展
这里是发展正文
第三章 结尾
这里是结尾正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
	assertChapter(t, chapters[2], 3, "结尾", "这里是结尾正文")
}

func TestParseChaptersLessThanThreeChapters(t *testing.T) {
	input := `第1章 开端
这里是开端正文
第二章 发展
这里是发展正文`

	chapters := ParseChapters(input)

	if len(chapters) != 2 {
		t.Fatalf("expected 2 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
}

func TestParseChaptersNoHeadingsReturnsEmptySlice(t *testing.T) {
	input := `前言
这里只是简介。
没有任何正式章节标题。`

	chapters := ParseChapters(input)

	if len(chapters) != 0 {
		t.Fatalf("expected empty chapter slice, got %d", len(chapters))
	}
}

func TestParseChaptersMergesPaginatedSameSection(t *testing.T) {
	input := `第一节：雨夜来信 (第1/2页)
第一页正文
第一节：雨夜来信 (第2/2页)
第二页正文`

	chapters := ParseChapters(input)

	if len(chapters) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "雨夜来信", "第一页正文\n第二页正文")
}

func TestParsePaginatedSections(t *testing.T) {
	input := `第四节：古月方源！ (第1/2页)
正文A

第四节：古月方源！ （第２／２页）
正文B

第五节：人祖三蛊，希望开窍 (第1/2页)
正文C

第五节：人祖三蛊，希望开窍 （第２／２页）
正文D

第六节：未来的路，会很精彩 (第1/2页)
正文E

第六节：未来的路，会很精彩 （第２／２页）
正文F`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapterContains(t, chapters[0], 1, "古月方源！", "正文A", "正文B")
	assertChapterContains(t, chapters[1], 2, "人祖三蛊，希望开窍", "正文C", "正文D")
	assertChapterContains(t, chapters[2], 3, "未来的路，会很精彩", "正文E", "正文F")
}

func TestParseInlinePaginatedSectionHeading(t *testing.T) {
	input := `第四节：古月方源！ (第1/2页) 朝阳升起来。
第四节：古月方源！ (第2/2页) 山雾不是很浓。
第五节：人祖三蛊，希望开窍 (第1/2页) 方源踏上对岸。
第五节：人祖三蛊，希望开窍 (第2/2页) 希望蛊汇入体内。
第六节：未来的路，会很精彩 (第1/2页) 方源测试结束。
第六节：未来的路，会很精彩 (第2/2页) 方正走到四十三步。`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapterContains(t, chapters[0], 1, "古月方源！", "朝阳升起来。", "山雾不是很浓。")
	assertChapterContains(t, chapters[1], 2, "人祖三蛊，希望开窍", "方源踏上对岸。", "希望蛊汇入体内。")
	assertChapterContains(t, chapters[2], 3, "未来的路，会很精彩", "方源测试结束。", "方正走到四十三步。")
}

func TestDoNotTreatStepsAsHeadings(t *testing.T) {
	input := `第1章 开端
方源走到第二十七步。
第二十七步
第二章 发展
古月赤城走到三十六步。
三十六步
第三章 结果
古月方正走到四十三步。
四十三步`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapterContains(t, chapters[0], 1, "开端", "第二十七步")
	assertChapterContains(t, chapters[1], 2, "发展", "三十六步")
	assertChapterContains(t, chapters[2], 3, "结果", "四十三步")
}

func TestParseMixedWhitespaceAndCRLF(t *testing.T) {
	input := "\uFEFF正文前言。\r\n　第四节：古月方源！ （第１／２页） 朝阳升起来。\r\n第四节：古月方源！ （第２／２页） 山雾不是很浓。\r\n正文末尾。　第五节：人祖三蛊，希望开窍 (第1/1页) 方源踏上对岸。\r第六节：未来的路，会很精彩 (第1/1页) 方正走到四十三步。"

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapterContains(t, chapters[0], 1, "古月方源！", "朝阳升起来。", "山雾不是很浓。", "正文末尾。")
	assertChapterContains(t, chapters[1], 2, "人祖三蛊，希望开窍", "方源踏上对岸。")
	assertChapterContains(t, chapters[2], 3, "未来的路，会很精彩", "方正走到四十三步。")
}

func TestParseChapterHeadingStableKeys(t *testing.T) {
	tests := []struct {
		line  string
		key   string
		title string
	}{
		{"第四节：古月方源！ (第1/2页)", "cn:四:节:古月方源！", "古月方源！"},
		{"第四节：古月方源！ （第２／２页）", "cn:四:节:古月方源！", "古月方源！"},
		{"第五节：人祖三蛊，希望开窍 (第1/2页)", "cn:五:节:人祖三蛊，希望开窍", "人祖三蛊，希望开窍"},
		{"Chapter 1: Beginning", "en:1:chapter:Beginning", "Beginning"},
	}

	for _, test := range tests {
		key, title, ok := parseChapterHeading(test.line)
		if !ok {
			t.Fatalf("expected heading %q to parse", test.line)
		}
		if key != test.key {
			t.Fatalf("expected key %q, got %q", test.key, key)
		}
		if title != test.title {
			t.Fatalf("expected title %q, got %q", test.title, title)
		}
	}
}

func assertChapter(t *testing.T, chapter Chapter, number int, title string, text string) {
	t.Helper()

	if chapter.Number != number {
		t.Fatalf("expected chapter number %d, got %d", number, chapter.Number)
	}
	if chapter.Title != title {
		t.Fatalf("expected chapter title %q, got %q", title, chapter.Title)
	}
	if chapter.Text != text {
		t.Fatalf("expected chapter text %q, got %q", text, chapter.Text)
	}
}

func assertChapterContains(t *testing.T, chapter Chapter, number int, title string, values ...string) {
	t.Helper()

	if chapter.Number != number {
		t.Fatalf("expected chapter number %d, got %d", number, chapter.Number)
	}
	if chapter.Title != title {
		t.Fatalf("expected chapter title %q, got %q", title, chapter.Title)
	}
	for _, value := range values {
		if !strings.Contains(chapter.Text, value) {
			t.Fatalf("expected chapter text %q to contain %q", chapter.Text, value)
		}
	}
}
