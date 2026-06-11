"use client";

import type { ReactNode } from "react";
import { useState } from "react";
import { Card, Segmented, Tag } from "antd";
import type { ChapterAnalysis, GenerateResponse, StoryBible } from "@/lib/api";

type ResultSectionsProps = {
  result: GenerateResponse | null;
  detailsOnly?: boolean;
};

type DetailSection = "Chapter Analysis" | "Factual Anchors" | "Story Bible";

export function ResultSections({ result, detailsOnly = false }: ResultSectionsProps) {
  const [section, setSection] = useState<DetailSection>("Chapter Analysis");

  if (!result || !detailsOnly) return null;

  return (
    <section className="detail-area">
      <div className="detail-heading">
        <div>
          <span className="section-kicker">03 / PIPELINE ARTIFACTS</span>
          <h2>查看中间产物</h2>
          <p>追踪从章节分析到全局故事资料的结构化推理结果。</p>
        </div>
        <Segmented<DetailSection>
          options={["Chapter Analysis", "Factual Anchors", "Story Bible"]}
          value={section}
          onChange={setSection}
        />
      </div>
      <div className="detail-content">
        {section === "Chapter Analysis" ? (
          <ChapterAnalyses analyses={result.chapter_analyses} />
        ) : section === "Factual Anchors" ? (
          <FactualAnchors analyses={result.chapter_analyses} />
        ) : (
          <StoryBibleView storyBible={result.story_bible} />
        )}
      </div>
    </section>
  );
}

export function OverviewContent({ result }: { result: GenerateResponse }) {
  const provider = result.meta?.ai_provider ?? "unknown";
  const issueCount = result.fidelity_result.issues.length;

  return (
    <>
      <div className="panel-title-row">
        <div>
          <span className="section-kicker">RUN SUMMARY</span>
          <h2>生成总览</h2>
        </div>
        <Tag className={provider === "real" ? "mode-badge real-mode" : "mode-badge mock-mode"}>
          <i />{provider === "real" ? "Real LLM" : provider === "mock" ? "Mock Mode" : "Unknown"}
        </Tag>
      </div>
      <div className="metric-grid">
        <Metric label="AI provider" value={provider} />
        <Metric label="AI model" value={result.meta?.ai_model || "未配置"} />
        <Metric label="章节数量" value={String(result.chapter_count).padStart(2, "0")} />
        <Metric label="Schema 校验" value={result.validation.passed ? "通过" : "失败"} status={result.validation.passed} />
        <Metric label="Fidelity Check" value={result.fidelity_result.passed ? "通过" : "有风险"} status={result.fidelity_result.passed} />
        <Metric label="Issues" value={String(issueCount).padStart(2, "0")} status={issueCount === 0} />
      </div>
      <div className="completion-panel">
        <div className="completion-copy">
          <span className="section-kicker">PIPELINE STATUS</span>
          <h3>工作流执行完成</h3>
          <p>章节事实已汇总，剧本结构与一致性检查结果可供审阅。</p>
        </div>
        <div className="completion-list">
          <p><i />Chapter Analysis</p>
          <p><i />Story Bible</p>
          <p><i />YAML Export</p>
        </div>
      </div>
    </>
  );
}

function Metric({ label, value, status }: { label: string; value: string | number; status?: boolean }) {
  return (
    <div className="metric">
      <span>{label}</span>
      <strong className={status === undefined ? "" : status ? "validation-ok" : "validation-bad"}>{value}</strong>
    </div>
  );
}

function ChapterAnalyses({ analyses }: { analyses: ChapterAnalysis[] }) {
  return (
    <Card className="tool-card artifact-card">
      <div className="card-heading">
        <span className="section-kicker">MAP / CHAPTER ANALYSIS</span>
        <h2>逐章结构化分析</h2>
        <p>每章单独分析人物、地点、事件、冲突和候选场景，用于保留长文本细节。</p>
      </div>
      <div className="stack">
        {analyses.map((chapter) => (
          <details className="chapter-card" key={`${chapter.chapter_number}-${chapter.chapter_title}`} open={chapter.chapter_number === 1}>
            <summary>
              <span>CH.{String(chapter.chapter_number).padStart(2, "0")}</span>
              <strong>{chapter.chapter_title}</strong>
              <small>{chapter.scene_candidates.length} 个候选场景</small>
            </summary>
            <div className="chapter-content">
              <p className="chapter-summary">{chapter.summary}</p>
              <FieldList title="角色分析">
                {chapter.characters.map((character) => (
                  <li key={character.name}>
                    <strong>{character.name}</strong>：{character.role_in_chapter}
                    {character.traits.length ? `；特征：${character.traits.join("、")}` : ""}
                    {character.state_change ? `；变化：${character.state_change}` : ""}
                  </li>
                ))}
              </FieldList>
              <div className="chapter-columns">
                <FieldList title="地点" items={chapter.locations} />
                <FieldList title="关键事件" items={chapter.key_events} />
                <FieldList title="冲突" items={chapter.conflicts} />
              </div>
              <div className="field-block anchor-inline">
                <h4>事实锚点</h4>
                {chapter.factual_anchors.length ? (
                  <ul>{chapter.factual_anchors.map((anchor) => <li key={anchor}>{anchor}</li>)}</ul>
                ) : (
                  <p className="muted">本章暂无事实锚点。</p>
                )}
              </div>
              <div className="field-block">
                <h4>候选场景</h4>
                <div className="mini-grid">
                  {chapter.scene_candidates.map((scene, index) => (
                    <div className="mini-card" key={`${scene.location}-${index}`}>
                      <div className="mini-card-top"><strong>{scene.location || "未指定地点"}</strong><span>{scene.time || "未指定时间"}</span></div>
                      <p>{scene.purpose}</p>
                      <small>角色：{scene.characters.join("、") || "无"}</small>
                      <small>事件：{scene.key_events.join("、") || "无"}</small>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </details>
        ))}
      </div>
    </Card>
  );
}

function FactualAnchors({ analyses }: { analyses: ChapterAnalysis[] }) {
  return (
    <Card className="tool-card artifact-card anchors-panel">
      <div className="card-heading">
        <span className="section-kicker">VERIFY / FACTUAL ANCHORS</span>
        <h2>事实锚点</h2>
        <p>记录原文中必须保留的硬事实，例如关键数字、人物关系、地点、事件结果和关键短句，用于约束后续剧本生成。</p>
      </div>
      <div className="stack compact-stack">
        {analyses.map((chapter) => (
          <article className="anchor-group" key={`${chapter.chapter_number}-${chapter.chapter_title}-anchors`}>
            <span>CH.{String(chapter.chapter_number).padStart(2, "0")}</span>
            <h3>{chapter.chapter_title}</h3>
            {chapter.factual_anchors.length ? (
              <ul>{chapter.factual_anchors.map((anchor) => <li key={anchor}>{anchor}</li>)}</ul>
            ) : (
              <p className="muted">本章暂无事实锚点。</p>
            )}
          </article>
        ))}
      </div>
    </Card>
  );
}

function StoryBibleView({ storyBible }: { storyBible: StoryBible }) {
  return (
    <Card className="tool-card artifact-card">
      <div className="card-heading">
        <span className="section-kicker">REDUCE / STORY BIBLE</span>
        <h2>全局故事资料</h2>
        <p>将多章分析结果合并为全局故事资料，统一人物、时间线、主线冲突和分场计划，减少多章节改编中的剧情漂移。</p>
      </div>
      <div className="story-header">
        <span className="section-kicker">LOGLINE</span>
        <h3>{storyBible.title}</h3>
        <p>{storyBible.logline}</p>
      </div>
      <FieldList title="全局角色">
        {storyBible.global_characters.map((character) => (
          <li key={character.id}><strong>{character.name}</strong>（{character.role}）：{character.motivation}</li>
        ))}
      </FieldList>
      <FieldList title="时间线">
        {storyBible.timeline.map((item) => <li key={`${item.chapter_number}-${item.event}`}>第 {item.chapter_number} 章：{item.event}</li>)}
      </FieldList>
      <div className="field-block">
        <h4>主线冲突</h4>
        <p>{storyBible.main_conflict}</p>
      </div>
      <div className="field-block">
        <h4>分场计划</h4>
        <div className="mini-grid">
          {storyBible.scene_plan.map((scene) => (
            <div className="mini-card" key={scene.id}>
              <div className="mini-card-top"><strong>{scene.id}</strong><span>CH.{String(scene.source_chapter).padStart(2, "0")}</span></div>
              <small>{scene.location} / {scene.time}</small>
              <p>{scene.summary}</p>
              <small>角色：{scene.characters.join("、") || "无"}</small>
            </div>
          ))}
        </div>
      </div>
    </Card>
  );
}

function FieldList({ title, items, children }: { title: string; items?: string[]; children?: ReactNode }) {
  const values = items ?? [];
  if (!children && values.length === 0) return null;
  return (
    <div className="field-block">
      <h4>{title}</h4>
      <ul>{children ?? values.map((item) => <li key={item}>{item}</li>)}</ul>
    </div>
  );
}
