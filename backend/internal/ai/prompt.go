package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/fidelity"
	"ai-showrunner-workbench/internal/novel"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"
)

const systemPrompt = "你是小说改编剧本的结构化分析助手。你必须只输出合法 JSON，不输出 markdown、解释、注释或 YAML。"

const faithfulAdaptationPrinciples = `1. 不得新增原文没有明确出现的具体事实。禁止新增原文没有的道具、动作、人物关系、地点、数量、因果解释、作弊手段、法术、物品或动机。
2. 可以做剧本化改编，但必须忠实原文事实。允许压缩情节、合并相近场景、把心理活动改写成可表演的低风险动作或短对白；禁止改变事件结果、人物立场、原文因果，禁止新增关键剧情，禁止把模糊暗示改成确定事实。
3. 保留原小说气质。尽量保留原文专有名词、人物姓名、地点名称、修行体系名词、关键设定名词、核心意象和整体气氛。
4. 经典短句可以少量保留，但不能整段照抄。对白优先参考原文已有说法；没有原文台词时可以生成改编型对白，但不能加入新事实。
5. 对不确定信息使用“似乎”“暗示”“存在嫌疑”“可能”“背后有隐情”等模糊表达，不要编造确定细节。

示例：原文只说“古月赤城作弊”或“方源知道古月赤城作弊信息”时，允许写“古月赤城的测试结果存在作弊嫌疑”或“方源知道古月赤城背后有隐情”；禁止写“古月赤城偷偷服下作弊丹”“古月赤城袖中藏着作弊蛊”“古月赤练当场施法帮助作弊”。`

func BuildChapterAnalysisPrompt(chapter novel.Chapter) string {
	return fmt.Sprintf(`请根据下面的单章小说文本，输出严格符合 analysis.ChapterAnalysis 的 JSON。

忠实改编原则：
%s

必须包含字段：
- chapter_number
- chapter_title
- summary
- characters
- locations
- key_events
- conflicts
- scene_candidates
- factual_anchors

characters 每一项必须包含：
- name
- role_in_chapter
- traits
- state_change

scene_candidates 每一项必须包含：
- location
- time
- purpose
- characters
- key_events

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 不要省略字段
- 字段名必须使用上述 snake_case JSON 字段名
- traits、locations、key_events、conflicts、characters 等数组字段必须始终输出 JSON 数组，即使只有一个元素也不能输出字符串
- 如果信息不足，也要给出合理的空数组或简短说明，但不要缺字段
- 只提取原文章节中明确出现或强烈暗示的信息
- characters、locations、key_events、conflicts、scene_candidates 都不得新增原文没有的事实
- scene_candidates 的 purpose 可以概括叙事功能，但不能新增剧情
- 对心理活动可以概括，但不能改写成原文没有发生的外部行为
- factual_anchors 只记录原文明确出现或强烈直接支持的硬事实，不要加入推测、猜测或画面补写
- factual_anchors 每条应是短句，适合给后续剧本生成做事实约束
- factual_anchors 优先记录关键数字、步数、资质等级、人物关系、地点、事件结果、明确出现的专有名词、原文明确出现的关键短句
- factual_anchors 不能写原文没有出现的细节

characters 示例：
[
  {
    "name": "角色名",
    "role_in_chapter": "本章作用",
    "traits": ["性格特征1", "性格特征2"],
    "state_change": "状态变化"
  }
]

章节编号：%d
章节标题：%s
章节正文：
%s`, faithfulAdaptationPrinciples, chapter.Number, chapter.Title, chapter.Text)
}

func BuildStoryBiblePrompt(analyses []analysis.ChapterAnalysis) string {
	payload, _ := json.MarshalIndent(analyses, "", "  ")
	return fmt.Sprintf(`请根据所有章节分析结果，合并成 story.StoryBible JSON。

忠实改编原则：
%s

必须输出字段：
- title
- logline
- global_characters
- timeline
- main_conflict
- scene_plan

global_characters 每一项必须包含：
- id
- name
- role
- motivation

timeline 每一项必须包含：
- chapter_number
- event

scene_plan 每一项必须包含：
- id
- source_chapter
- summary
- location
- time
- characters

重点：
- 统一人物名称，避免同一人物在不同章节漂移
- 为每个主要角色生成稳定 id，建议使用 char_ 开头的小写英文/拼音/下划线
- 合并时间线
- 提炼主线冲突
- 生成适合剧本改编的 title 和 logline
- 生成至少 3 个适合剧本改编的 scene_plan
- 合并多章分析时，不要发明新主线
- title / logline 可以概括，但不能改变故事类型和核心冲突
- global_characters 只能来自章节分析中出现过的人物
- timeline 只能基于已分析章节的事件
- main_conflict 可以抽象概括，但不能新增反派计划或隐藏设定
- scene_plan 只能来自 scene_candidates 和 key_events

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 字段名必须符合上述 JSON 字段名
- 不要缺字段

章节分析 JSON：
%s`, faithfulAdaptationPrinciples, string(payload))
}

type screenplayPromptPayload struct {
	StoryBible      story.StoryBible       `json:"story_bible"`
	ChapterAnalyses []screenplayFactAnchor `json:"chapter_analyses"`
}

type screenplayFactAnchor struct {
	ChapterNumber   int                       `json:"chapter_number"`
	ChapterTitle    string                    `json:"chapter_title"`
	Summary         string                    `json:"summary"`
	FactualAnchors  []string                  `json:"factual_anchors"`
	KeyEvents       []string                  `json:"key_events"`
	Conflicts       []string                  `json:"conflicts"`
	SceneCandidates []analysis.SceneCandidate `json:"scene_candidates"`
}

func BuildScreenplayPrompt(bible story.StoryBible, analyses []analysis.ChapterAnalysis) string {
	payload, _ := json.MarshalIndent(screenplayPromptPayload{
		StoryBible:      bible,
		ChapterAnalyses: buildScreenplayFactAnchors(analyses),
	}, "", "  ")
	return fmt.Sprintf(`请根据 StoryBible 和 ChapterAnalyses 生成 screenplay.Screenplay JSON。

这是忠实于原文事实的剧本化改编初稿，不是完全二创。

ChapterAnalyses 是最终剧本生成的事实锚点。生成 scenes、dialogues 和 actions 时，必须优先遵守每章的 factual_anchors、summary、key_events、conflicts、scene_candidates。

忠实改编原则：
%s

必须输出字段：
- title
- source_chapters
- characters
- scenes

source_chapters 每一项必须包含：
- number
- title
- summary

characters 每一项必须包含：
- id
- name
- role
- description

scenes 每一项必须包含：
- id
- source_chapter
- location
- time
- summary
- characters
- dialogues
- actions

dialogues 每一项必须包含：
- character
- emotion
- line

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- characters 字段里引用角色 id
- dialogue.character 也使用角色 id，并尽量和 characters 中的 id 一致
- 至少生成 3 个 scenes
- 每个 scene 至少有 1 条 dialogue
- 不要缺少 validator 要求的字段
- source_chapters.title 不要重新命名；如果输入中已有章节标题，应原样保留
- 可以生成改编对白，但不得加入新事实
- dialogues.line 如果来自原文经典短句，可以短句保留；如果是改编对白，要符合人物处境和原文气质
- actions 不得新增原文没有的具体道具和行为
- 如果原文只是心理活动，actions 可以写为“沉默、观察、停顿、转身”等低风险表演动作，不要编造新动作
- 不要新增原文没有出现的人物

关键事实保留规则：
- 所有具体事实必须可追溯到 StoryBible、ChapterAnalysis.factual_anchors、ChapterAnalysis.key_events 或 ChapterAnalysis.scene_candidates
- 以下内容不得改写：人物姓名、人物关系、地点、事件结果、资质等级、步数、数量、章节归属、明确出现过的专有名词
- 如果没有依据，只能概括表达，不能补具体细节
- “二十七步”不能改成“不足十步”
- “三十六步”不能改成“三十步”
- “四十三步”不能改成“四十步”
- “丙等资质”不能改成“丁等资质”
- “甲等资质”不能改成“乙等资质”
- “希望蛊数量偏少”不能改成“只出现一只希望蛊”
- 不得把“方正”写成“古月漠尘的孙儿”，除非 factual_anchors 明确支持
- 不得新增“通过蛊器观察”“袖口颤动”“第十七步额头冒汗”等无依据观察方式、身体反应或细节
- 不得新增“作弊丹、作弊蛊、符纸”等具体作弊道具

章节归属规则：
- 每个 scene.source_chapter 必须和 ChapterAnalyses 中的事实来源一致
- 不要把后续章节的宣布结果、评价或反应提前塞进前一章 scene
- 第 2 章如果包含“方源走到二十七步、希望蛊聚集、体内开窍成功、外界认为希望蛊数量偏少”，就只能把这些写入第 2 章相关 scene
- 第 3 章如果包含“方源测试结果引发失望、方正测试甲等资质、家老争夺培养权”，这些宣布结果和争夺反应不要提前写进第 2 章 scene

改编台词规则：
- 可以生成剧本化改编台词
- 如果是原文经典短句，可以短句保留，例如“未来的路，会很精彩呢。”
- 如果不是原文原句，台词必须只表达原文已经存在的情绪、关系或冲突
- 改编台词不得加入新事实，不要伪造原文没有说过的关键结论
- “别怕，只管走。”这类台词只能作为情绪性改编，不能承载新事实
- 禁止写“我已经服下作弊丹”这类原文没有的具体事实

actions 规则：
- 每条 action 必须是完整可拍摄动作句
- 不要拆成过短碎片
- 不要单独写“若有所思”“目光微凝”这类碎片，除非和完整动作合并
- 每条 action 建议 15 到 40 个中文字符
- actions 应优先来自原文叙述或对心理活动的低风险外化
- 心理活动可以外化为“停步、沉默、观察、转身、低头”等动作，但不要编造新行为
- 坏例子：["方源站在测试队列中", "若有所思"]
- 好例子：["方源站在测试队列中，平静地注视古月赤城走入花海。", "听到乙等资质的宣布后，方源没有出声，只是收回目光。"]

StoryBible 和 ChapterAnalyses JSON：
%s`, faithfulAdaptationPrinciples, string(payload))
}

func buildScreenplayFactAnchors(analyses []analysis.ChapterAnalysis) []screenplayFactAnchor {
	anchors := make([]screenplayFactAnchor, 0, len(analyses))
	for _, item := range analyses {
		anchors = append(anchors, screenplayFactAnchor{
			ChapterNumber:   item.ChapterNumber,
			ChapterTitle:    item.ChapterTitle,
			Summary:         item.Summary,
			FactualAnchors:  item.FactualAnchors,
			KeyEvents:       item.KeyEvents,
			Conflicts:       item.Conflicts,
			SceneCandidates: item.SceneCandidates,
		})
	}
	return anchors
}

type fidelityPromptPayload struct {
	Screenplay      screenplay.Screenplay    `json:"screenplay"`
	StoryBible      story.StoryBible         `json:"story_bible"`
	ChapterAnalyses []screenplayFactAnchor   `json:"chapter_analyses"`
	FidelityIssues  []fidelity.FidelityIssue `json:"fidelity_issues,omitempty"`
}

func BuildFidelityCheckPrompt(current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis) string {
	payload, _ := json.MarshalIndent(fidelityPromptPayload{
		Screenplay:      current,
		StoryBible:      bible,
		ChapterAnalyses: buildScreenplayFactAnchors(analyses),
	}, "", "  ")

	return fmt.Sprintf(`请检查 screenplay 是否忠实于 StoryBible 和 ChapterAnalyses 中的事实锚点，输出 fidelity.FidelityResult JSON。

只检查事实一致性，不评价文笔。

重点检查：
1. 关键数字是否改错
2. 资质等级是否改错
3. 人物关系是否改错
4. 地点是否凭空新增
5. 道具是否凭空新增
6. 观察方式是否凭空新增
7. 身体反应是否凭空新增
8. 章节归属是否混乱
9. actions 是否补写原文没有的具体事实
10. dialogues 是否承载了原文没有的关键信息
11. 如果 screenplay 的 title、characters 或 scenes 为空，必须报告 high severity issue，不能判定 passed=true

事实依据只能来自：
- StoryBible
- ChapterAnalysis.factual_anchors
- ChapterAnalysis.key_events
- ChapterAnalysis.scene_candidates

unsupported claims 示例：
- 把“二十七步”改成其他数字
- 把“方正”写成“古月漠尘的孙儿”
- 新增“通过蛊器观察”
- 新增“袖口颤动”
- 新增“第十七步额头冒汗”
- 新增“作弊丹、作弊蛊、符纸”等具体作弊道具

返回格式：
{
  "passed": true,
  "issues": []
}

或：
{
  "passed": false,
  "issues": [
    {
      "field": "scenes[2].dialogues[4].line",
      "severity": "high",
      "problem": "问题说明",
      "suggestion": "修复建议"
    }
  ]
}

severity 只能是 low、medium、high。

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 不要省略字段

输入 JSON：
%s`, string(payload))
}

func BuildFidelityRepairPrompt(current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis, result fidelity.FidelityResult) string {
	payload, _ := json.MarshalIndent(fidelityPromptPayload{
		Screenplay:      current,
		StoryBible:      bible,
		ChapterAnalyses: buildScreenplayFactAnchors(analyses),
		FidelityIssues:  result.Issues,
	}, "", "  ")

	return fmt.Sprintf(`请根据 fidelity_issues 修复 screenplay.Screenplay JSON。

修复规则：
- 只修复 issues 指出的事实问题
- 不改变 YAML 主 Schema
- 不新增无依据事实
- 不删除必填字段
- 返回完整 Screenplay JSON
- 缺失或错误事实只能从 StoryBible、ChapterAnalysis.factual_anchors、key_events、scene_candidates 中补
- dialogues 可以保守改写，但不得承载新事实
- actions 必须是完整可拍摄动作句，且不得新增具体道具、观察方式、身体反应或亲属关系
- source_chapters 可保留当前值，后端会再次用 ChapterAnalyses 确定性覆盖

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 不要省略 validator 必需字段

输入 JSON：
%s`, string(payload))
}

func BuildRepairJSONPrompt(originalPrompt string, rawOutput string, parseError error, schema string) string {
	return fmt.Sprintf(`你上一次返回的内容无法解析为目标 Go struct JSON。请基于原始任务、解析错误和目标字段要求，返回修复后的合法 JSON。

要求：
- 只返回修复后的 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 字段名必须符合 json tag
- 不要省略必填字段
- 数组字段必须输出 JSON 数组，不能输出字符串
- 修复时不得改变原始任务中的事实
- 如果是修复 screenplay JSON，缺失字段必须从原始任务里的 StoryBible 和 ChapterAnalyses 事实锚点补充
- 不得为了通过解析而新增原文没有的道具、数量、事件或人物

目标字段要求：
%s

json.Unmarshal 错误：
%s

原始任务：
%s

模型第一次返回内容：
%s`, schema, parseError.Error(), truncateText(originalPrompt, repairPromptLimit), truncateText(rawOutput, repairPromptLimit))
}

func BuildRepairScreenplayPrompt(originalPrompt string, current screenplay.Screenplay, validationErrors []string) string {
	payload, _ := json.MarshalIndent(current, "", "  ")
	return fmt.Sprintf(`你上一次返回的 screenplay.Screenplay JSON 已经可以解析，但没有通过后端 Validate。请基于原始任务、当前 JSON 和校验错误，返回修复后的 screenplay.Screenplay JSON。

要求：
- 只返回修复后的 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 字段名必须符合 json tag
- 不要省略必填字段
- 不要删除 validator 需要的字段
- source_chapters 必须非空，每项 number > 0 且 title 非空
- characters 必须非空，每项 id、name、role 非空
- scenes 必须非空，每项 id、location、time、summary 非空
- 每个 scene.characters 必须非空
- 每个 scene.dialogues 必须非空
- 每条 dialogue.character 和 dialogue.line 必须非空
- 修复时不能为了通过 validator 而编造字段内容
- 缺字段应从原始任务里的 StoryBible 和 ChapterAnalyses 事实锚点补充
- 不得引入新的道具、数量、事件、人物关系、地点或章节归属
- 关键数字、资质等级、人物关系、地点和事件结果必须保持原样
- actions 仍必须是完整可拍摄动作句，不要拆成短碎片

忠实改编原则：
%s

校验错误：
%s

原始任务：
%s

当前 screenplay JSON：
%s`, faithfulAdaptationPrinciples, strings.Join(validationErrors, "\n"), truncateText(originalPrompt, repairPromptLimit), truncateText(string(payload), repairPromptLimit))
}

func chapterAnalysisSchemaDescription() string {
	return `analysis.ChapterAnalysis:
{
  "chapter_number": number,
  "chapter_title": string,
  "summary": string,
  "characters": [
    {
      "name": string,
      "role_in_chapter": string,
      "traits": [string],
      "state_change": string
    }
  ],
  "locations": [string],
  "key_events": [string],
  "conflicts": [string],
  "factual_anchors": [string],
  "scene_candidates": [
    {
      "location": string,
      "time": string,
      "purpose": string,
      "characters": [string],
      "key_events": [string]
    }
  ]
}`
}

func storyBibleSchemaDescription() string {
	return `story.StoryBible:
{
  "title": string,
  "logline": string,
  "global_characters": [
    {
      "id": string,
      "name": string,
      "role": string,
      "motivation": string
    }
  ],
  "timeline": [
    {
      "chapter_number": number,
      "event": string
    }
  ],
  "main_conflict": string,
  "scene_plan": [
    {
      "id": string,
      "source_chapter": number,
      "summary": string,
      "location": string,
      "time": string,
      "characters": [string]
    }
  ]
}`
}

func screenplaySchemaDescription() string {
	return `screenplay.Screenplay:
{
  "title": string,
  "source_chapters": [
    {
      "number": number,
      "title": string,
      "summary": string
    }
  ],
  "characters": [
    {
      "id": string,
      "name": string,
      "role": string,
      "description": string
    }
  ],
  "scenes": [
    {
      "id": string,
      "source_chapter": number,
      "location": string,
      "time": string,
      "summary": string,
      "characters": [string],
      "dialogues": [
        {
          "character": string,
          "emotion": string,
          "line": string
        }
      ],
      "actions": [string]
    }
  ]
}`
}

func fidelityResultSchemaDescription() string {
	return `fidelity.FidelityResult:
{
  "passed": boolean,
  "issues": [
    {
      "field": string,
      "severity": "low" | "medium" | "high",
      "problem": string,
      "suggestion": string
    }
  ]
}`
}
