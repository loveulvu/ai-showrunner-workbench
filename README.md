# Novel to Screenplay AI

一个基于多阶段 AI Workflow 的小说转结构化 YAML 剧本工具。

## Demo Video

- 在线演示视频（Bilibili）：[Watch Demo Video](https://www.bilibili.com/video/BV1rGEb6aE7u/)
- 备用下载链接（GitHub Release）：[Download Demo Video](https://github.com/loveulvu/novel-to-screenplay-ai/releases/tag/v0.1.0)

Demo 展示完整流程：

```text
Novel Text → Chapter Analysis → Factual Anchors → Story Bible → Fidelity Check → YAML Screenplay Export
```

## Project Background

本项目对应 **七牛云 × XEngineer 暑期实训营第三批次议题三：AI 小说转剧本工具**。

题目要求输入三个章节以上的小说文本，自动转换为结构化剧本 YAML，并定义 YAML Schema 及说明其设计原因。

## Core Features

- 多章节小说文本输入
- 支持中文“第 X 章 / 第 X 节”和英文 `Chapter` 解析
- 支持分页章节合并
- 逐章 LLM 结构化分析
- Story Bible 全局故事合并
- Factual Anchors 事实锚点
- Fidelity Check 事实一致性检查与一次定向 Repair
- YAML Schema 校验
- YAML 预览、复制和下载
- Mock / Real AI Provider 切换

## Why It Is Not Just a Prompt Wrapper

本项目不是一次性把小说全文交给大模型直接生成 YAML，而是采用固定、可解释的多阶段 AI Workflow：

1. Parse Chapters
2. Analyze Each Chapter
3. Extract Factual Anchors
4. Merge Story Bible
5. Generate Screenplay JSON
6. Fidelity Check & Repair
7. Schema Validate
8. Export YAML

采用该设计的原因：

- 长文本一次生成容易遗漏早期章节细节
- 多章节改编容易出现人物名称、关系和设定漂移
- LLM 可能产生无原文依据的事实幻觉
- YAML 输出需要经过结构校验，才能稳定地被程序读取
- Factual Anchors 和 Fidelity Check 用于降低事实错误风险

逐章分析对应 Map-Reduce 中的 **Map** 阶段；Story Bible 合并对应 **Reduce** 阶段。每个阶段都有明确输入输出，便于展示、调试和替换模型。

## Architecture

```text
Novel Text
    ↓
Parse Chapters
    ↓
Analyze Each Chapter
    ↓
Extract Factual Anchors
    ↓
Merge Story Bible
    ↓
Generate Screenplay JSON
    ↓
Fidelity Check & Repair
    ↓
Schema Validate
    ↓
Export YAML
```

质量控制分为两部分：

- `Fidelity Check` 检查事实一致性，例如无依据事实、人物关系错误、关键数字错误和章节归属混乱。
- `Schema Validate` 检查结构完整性，确保最终剧本字段完整且可被程序读取。

## YAML Schema

最终 YAML 顶层结构：

- `title`
- `source_chapters`
- `characters`
- `scenes`

每个 `scene` 包含：

- `id`
- `source_chapter`
- `location`
- `time`
- `summary`
- `characters`
- `dialogues`
- `actions`

`source_chapters` 保留剧本与原小说章节的对应关系，`characters` 统一跨章节人物引用，`scenes` 承载可编辑的分场剧本。

完整字段定义与设计原因见 [`docs/schema.md`](docs/schema.md)。

## Tech Stack

- Frontend: Next.js + TypeScript
- Backend: Go + Gin
- LLM: OpenAI-compatible Chat Completions API
- Output: YAML
- Validation: Go struct + custom validator
- Demo: GitHub Release video

## Local Development

### Backend

```bash
cd backend
go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

默认地址：

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

## Environment Variables

项目支持 mock、real 和 qwen 三种 provider。参考 `.env.example` 创建本地 `.env`：

```env
AI_PROVIDER=qwen
AI_API_KEY=your_api_key_here
AI_BASE_URL=https://dashscope-intl.aliyuncs.com/compatible-mode/v1
AI_MODEL=qwen-plus
AI_TIMEOUT_SECONDS=180
```

### Qwen Cloud

Set `AI_PROVIDER=qwen` to use Qwen Cloud through its OpenAI-compatible Chat Completions endpoint. `AI_API_KEY` is required. When `AI_BASE_URL` or `AI_MODEL` is omitted, they default to `https://dashscope-intl.aliyuncs.com/compatible-mode/v1` and `qwen-plus`.

The backend sends requests to `{AI_BASE_URL}/chat/completions`. Keep the API key only in your local `.env`; `.env` files are ignored by Git and must not be committed.

To keep using the built-in deterministic provider, set `AI_PROVIDER=mock`. The generic `AI_PROVIDER=real` OpenAI-compatible provider remains available and still requires an explicit API key, base URL, and model.

- `AI_PROVIDER=mock`：使用内置 mock client，无需真实 API Key，适合开发和稳定 Demo。
- `AI_PROVIDER=real`：调用真实 OpenAI-compatible LLM。
- `AI_PROVIDER=qwen`：调用 Qwen Cloud OpenAI-compatible Chat Completions API。
- `AI_TIMEOUT_SECONDS`：单次真实 LLM HTTP 请求超时时间，默认 `180` 秒。
- `.env` 仅用于本地配置，不要提交到 GitHub。

### Environment Loading and PowerShell

At startup, the Go backend tries to load `.env` and `../.env` relative to its current working directory. Missing files are ignored, and startup logs show which files were loaded. Values already set in the process environment take precedence over `.env` values.

A `.env` file is read only by the Go backend when the application loads it. It does not automatically set variables in your PowerShell session, so direct `curl.exe` tests need PowerShell environment variables:

```powershell
$env:AI_API_KEY = "your_api_key_here"
$env:AI_BASE_URL = "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"

curl.exe "$env:AI_BASE_URL/chat/completions" `
  -H "Authorization: Bearer $env:AI_API_KEY" `
  -H "Content-Type: application/json" `
  -d '{"model":"qwen-plus","messages":[{"role":"user","content":"Hello"}]}'
```

Depending on the local network, Qwen Cloud access may require proxy variables for both the Go backend and direct PowerShell tests:

```powershell
$env:HTTP_PROXY = "http://127.0.0.1:7890"
$env:HTTPS_PROXY = "http://127.0.0.1:7890"
```

Startup diagnostics report only whether `AI_API_KEY`, `HTTP_PROXY`, and `HTTPS_PROXY` are set. Their values are never logged.

## API

### Health Check

```http
GET /api/health
```

### Generate Screenplay

```http
POST /api/generate
Content-Type: application/json
```

请求：

```json
{
  "novel_text": "..."
}
```

返回包含：

- `chapter_count`
- `chapter_analyses`
- `story_bible`
- `screenplay_json`
- `screenplay_yaml`
- `validation`
- `fidelity_result`
- `meta`

## Current Limitations

- 当前推荐输入 3～5 章
- 生成质量仍受 LLM 能力影响
- Fidelity Check 可降低风险，但不能保证 100% 无幻觉
- 暂不支持 PDF / Word 解析
- 暂不支持登录、历史记录和多人协作
- 当前 MVP 聚焦无登录的单次生成流程

## Repository Safety

- Demo 视频通过 GitHub Release 提供，不直接提交 MP4 到仓库
- API Key 通过本地 `.env` 配置，不进入 Git
- 仓库只保留 `.env.example` 占位配置
- README、docs、examples 和源码中不应包含真实 API Key
