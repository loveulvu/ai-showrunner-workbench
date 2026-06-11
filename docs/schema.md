# YAML Schema 设计说明

最终剧本使用结构化 YAML 表达。设计目标是：既保留小说来源和事实约束，又让结果适合阅读、编辑、复制、下载与程序校验。

## 顶层结构

```yaml
title: "雨夜旧剧院"
source_chapters:
  - number: 1
    title: "雨夜钥匙"
    summary: "林澈收到父亲留下的铜钥匙，线索指向海棠剧院。"
characters:
  - id: "char_lin_che"
    name: "林澈"
    role: "主角"
    description: "追查父亲失踪真相的青年。"
scenes:
  - id: "scene_001"
    source_chapter: 1
    location: "老街雨巷"
    time: "夜晚"
    summary: "林澈收到父亲留下的铜钥匙。"
    characters:
      - "char_lin_che"
    dialogues:
      - character: "char_lin_che"
        emotion: "迟疑"
        line: "这把钥匙像是在等我回来。"
    actions:
      - "林澈撑伞站在废弃邮箱前，铜钥匙滑入掌心。"
```

## 为什么有 `source_chapters`

`source_chapters` 用于保留剧本与原小说章节的对应关系。每个来源章节记录编号、原始标题和摘要，方便作者回溯事实来源、定位修改位置，也便于检查场景是否放错章节。

## 为什么有 `characters`

`characters` 是全局角色表，使用稳定 `id` 统一多章节人物引用，降低人物重名、别名、改名或名称漂移造成的混乱。场景和对白通过角色 `id` 引用同一个人物。

## 为什么有 `scenes`

剧本以场景为核心单位。每个场景包含来源章节、地点、时间、摘要、参与角色、对白和动作，使小说叙述能够转换为可展示、可编辑的分场剧本。

## 为什么 Dialogue 拆成 `character` / `emotion` / `line`

- `character`：明确说话人，支持角色引用校验。
- `emotion`：明确语气和表演方向。
- `line`：保存实际台词内容。

拆分后更适合程序读取、表演提示、后续格式转换和人工编辑。

## 为什么有 `actions`

`actions` 用于把小说叙述和心理活动转成可拍摄、可表演的动作描述。动作应以原文事实为依据，使用完整动作句，不应为了画面感新增无依据的道具、身体反应或关键行为。

## 字段约束

- `title`：剧本标题，必须非空。
- `source_chapters`：来源章节列表，必须非空。
- `source_chapters.number/title/summary`：来源编号、原始标题和章节贡献摘要。
- `characters`：全局角色列表，必须非空。
- `characters.id/name/role`：稳定标识、显示名称与角色功能，必须非空。
- `characters.description`：人物状态、动机或改编重点。
- `scenes`：分场列表，必须非空。
- `scene.id/location/time/summary`：场景核心字段，必须非空。
- `scene.source_chapter`：场景主要事实来源章节。
- `scene.characters`：参与场景的角色 `id` 列表。
- `scene.dialogues`：结构化对白列表。
- `scene.actions`：可拍摄、可表演的动作列表。

## 中间态为什么使用 JSON，最终为什么输出 YAML

中间态使用 JSON / Go struct，便于程序序列化、API 传输、类型约束和自定义校验。最终输出 YAML，是因为 YAML 层级清晰、可读性更好，更适合作者复制、下载、编辑和继续二次创作。

## ChapterAnalysis 中间态

逐章分析结果包含：

- `chapter_number` / `chapter_title` / `summary`
- `characters`
- `locations`
- `key_events`
- `conflicts`
- `scene_candidates`
- `factual_anchors`

这层中间态用于保留长文本细节，避免直接从全文跳到最终剧本。

## Factual Anchors

`factual_anchors` 记录原文中必须保留的硬事实，例如：

- 关键数字和数量
- 人物关系
- 地点和道具
- 事件结果
- 章节归属
- 专有名词和关键短句

这些事实用于约束生成和支撑 Fidelity Check。

## 双重质量检查

- `Schema Validate` 负责结构完整性，检查必需字段是否存在且可读取。
- `Fidelity Check` 负责事实一致性，检查是否出现无依据补写、人物关系错误、关键数字错误或章节归属混乱。

两者关注不同问题，不能互相替代。
