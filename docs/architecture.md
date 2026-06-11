# 系统架构与设计取舍

本项目采用固定、可控的 AI Workflow，将多章节小说转换为经过质量门控的结构化 YAML 剧本。

## 主流程

```text
Novel Text
→ Parse Chapters
→ Analyze Each Chapter
→ Extract Factual Anchors
→ Merge Story Bible
→ Generate Screenplay JSON
→ Fidelity Check & Repair
→ Schema Validate
→ Export YAML
```

代码主链路：

```text
ParseChapters
→ AnalyzeChapters
→ MergeStoryBible
→ GenerateScreenplay
→ CheckAndRepairOnce
→ ValidateScreenplay
→ ToYAML
```

## 固定 AI Workflow，不是复杂 Agent

本项目采用可控流水线，而不是让 LLM 自主决定工具、步骤和循环次数。比赛与 Demo 场景更重视稳定性、可解释性和可调试性；固定 Workflow 能明确展示每一步的输入、输出和失败位置，也便于单独替换模型或模块。

## 为什么不直接全文一次生成 YAML

长文本一次生成容易出现：

- 早期章节细节遗漏
- 跨章节人物名称或关系漂移
- 关键数字、地点和事件结果被改写
- 输出过长导致截断
- YAML 字段缺失或格式错误

因此系统先拆分章节和保存中间态，再生成最终剧本。

## Map-Reduce 思路

### Map：逐章分析

`Analyze Each Chapter` 对每章单独分析，提取章节摘要、人物、地点、关键事件、冲突、候选场景和事实锚点，生成 `ChapterAnalysis`。

独立分析能够保留长文本细节，也让每章结果可单独检查。

### Reduce：合并 Story Bible

`Merge Story Bible` 合并所有章节分析，统一跨章节人物、时间线、主线冲突和分场计划，生成全局 `StoryBible`。

Story Bible 为最终剧本提供稳定的全局视角，减少多章节改编中的剧情漂移。

## 模块职责

### Parse

位置：`backend/internal/novel`

识别中文章、节和英文 `Chapter` 标题，切分用户输入的多章节小说。输入少于三个章节时，API 返回明确错误。

### Analyze

位置：`backend/internal/analysis`

逐章调用 LLM，生成 `ChapterAnalysis` 和 `factual_anchors`。该阶段不直接生成最终剧本。

### Merge

位置：`backend/internal/story`

将多章分析结果合并为 `StoryBible`，统一全局人物、时间线、主冲突和场景计划。

### Generate

位置：`backend/internal/screenplay`

根据 Story Bible、Chapter Analysis 和事实锚点生成结构化 `Screenplay` JSON。

### Export

位置：`backend/internal/screenplay/yaml.go`

将内部 JSON / Go struct 转换为最终 YAML，供前端预览、复制和下载。

## 质量门控

### Fidelity Check

位置：`backend/internal/fidelity`

负责事实一致性，检查生成剧本是否出现：

- 无依据事实
- 人物关系错误
- 关键数字错误
- 地点、道具或事件结果错误
- 章节归属混乱

当存在中高风险问题时，系统只执行一次定向 Fidelity Repair，再重新检查。限制修复次数可以避免无限循环，并保留可解释的最终检查结果。

### Schema Validate

位置：`backend/internal/screenplay/validator.go`

负责结构完整性，检查标题、来源章节、角色、场景、对白等关键字段是否存在且可被程序读取。

Schema Validate 与 Fidelity Check 是两道不同质量门控：前者检查“结构是否正确”，后者检查“内容是否忠实”。

## Mock 与 Real Provider

后端通过 `AI_PROVIDER` 切换：

- `mock`：无需 API Key，适合开发、测试和稳定 Demo。
- `real`：调用 OpenAI-compatible Chat Completions API，执行真实逐章分析、合并、生成和事实检查。

两种模式都保留相同 Workflow 和 Schema Validate。

## 设计边界

当前版本聚焦于可解释的小说转剧本核心流程，不包含数据库、登录、历史记录、文件上传、流式输出、多人协作或桌面端。
