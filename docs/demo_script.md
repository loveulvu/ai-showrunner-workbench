# 3～5 分钟 Demo 录制脚本

## 0:00～0:25 项目一句话介绍

**画面：** 打开前端首页，展示标题、Workflow Stepper 和小说输入区。

**讲解：**

> 这是 AI Showrunner Workbench，一个基于多阶段 AI Workflow 的小说转结构化 YAML 剧本工具，对应七牛云 × XEngineer 暑期实训营第三批次议题三。

## 0:25～0:45 说明题目与设计思路

**画面：** 指向顶部 Workflow Stepper。

**讲解：**

> 本项目不是一次性把全文丢给模型生成 YAML，而是采用多阶段 AI Workflow。系统先做章节级结构化分析，再合并为 Story Bible，之后根据事实锚点生成剧本，并通过 Fidelity Check 和 Schema Validate 做质量控制。

## 0:45～1:10 展示输入小说

**操作：**

1. 点击“填入示例小说”。
2. 简要滚动输入框，展示三个章节标题。

**讲解：**

> 输入至少三个章节后，系统会先识别章节边界。每章会单独进入分析阶段，因此长文本不会直接挤进一次生成请求。

## 1:10～1:30 点击生成

**操作：** 点击“生成剧本 YAML”，展示 Stepper 和 Loading Pipeline 的阶段变化。

**讲解：**

> 生成过程依次完成章节分析、Story Bible 合并、事实检查和 YAML 导出。固定流程让每一步都可解释，也方便定位问题。

## 1:30～2:05 解释 Chapter Analysis

**操作：** 在 `Chapter Analysis` 页签展开第一章。

**讲解：**

> Chapter Analysis 是 Map 阶段。系统每章单独提取人物、地点、事件、冲突和候选场景，用于保留长文本细节。这里还能看到角色状态变化和场景改编建议。

## 2:05～2:30 解释 Factual Anchors

**操作：** 切换到 `Factual Anchors`。

**讲解：**

> 事实锚点记录原文中必须保留的硬事实，例如关键数字、人物关系、地点、事件结果和关键短句。后续剧本生成和 Fidelity Check 都会使用这些锚点，减少模型改错事实。

## 2:30～2:55 解释 Story Bible

**操作：** 切换到 `Story Bible`。

**讲解：**

> Story Bible 是 Reduce 阶段。它把多章分析合并为全局故事资料，统一人物、时间线、主线冲突和分场计划，减少多章节改编中的剧情漂移。

## 2:55～3:25 解释 Fidelity Check 与 Schema Validate

**操作：** 切换结果区的 `Quality` 页签。

**讲解：**

> 这里有两道质量门控。Fidelity Check 检查无依据事实、人物关系错误、关键数字错误和章节归属混乱；Schema Validate 检查最终结构是否完整、是否能被程序读取。两者关注的问题不同，不能互相替代。

## 3:25～3:55 展示 YAML 输出

**操作：** 切换到 `YAML` 页签，滚动展示 `source_chapters`、`characters` 和 `scenes`。

**讲解：**

> 最终剧本使用 YAML 输出。`source_chapters` 保留与原小说章节的对应关系，`characters` 统一人物引用，`scenes` 保存地点、时间、对白和动作。中间态使用 JSON 和 Go struct 方便校验，YAML 更适合阅读和编辑。

## 3:55～4:10 复制与下载

**操作：**

1. 点击“复制 YAML”。
2. 点击“下载 screenplay.yaml”。

**讲解：**

> 结果可以直接复制或下载，方便作者继续编辑和二次创作。

## 4:10～4:35 总结项目亮点

**画面：** 返回 Overview 或停留在完整结果页。

**讲解：**

> 项目的核心亮点不是单次 Prompt，而是可解释的长文本改编流程：章节级 Map 分析、Story Bible Reduce 合并、事实锚点约束、Fidelity Check 事实检查和 Schema Validate 结构校验。当前版本同时支持 mock 与 real provider，适合稳定演示，也可以接入真实 OpenAI-compatible 模型。
