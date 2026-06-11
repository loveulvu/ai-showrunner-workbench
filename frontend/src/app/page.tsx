"use client";

import { useEffect, useState } from "react";
import { Alert, ConfigProvider, Tag } from "antd";
import { GeneratePanel } from "@/components/GeneratePanel";
import { NovelInput } from "@/components/NovelInput";
import { OutputPanel } from "@/components/OutputPanel";
import { ResultSections } from "@/components/ResultSections";
import { WorkflowStepper } from "@/components/WorkflowStepper";
import { generateScreenplay } from "@/lib/api";
import type { GenerateResponse } from "@/lib/api";

const sampleText = `第一章 雨夜钥匙
林澈在雨夜回到老街。巷口的邮箱早已停用，却在这天晚上吐出一封没有署名的信。
信封里只有一把铜钥匙和半张旧剧票。剧票背面写着父亲熟悉的字迹：海棠剧院，午夜之后。

第二章 旧剧院
许岚不放心林澈独自前往，带着手电陪他走进废弃的海棠剧院。
他们在舞台下方找到一台仍在运转的旧时钟，齿轮之间夹着另一半剧票。

Chapter 3 舞台对峙
顾衡在灯光亮起时出现，要求林澈交出铜钥匙。
林澈握紧钥匙，决定启动时钟，查清父亲消失的真相。`;

export default function Home() {
  const [novelText, setNovelText] = useState("");
  const [result, setResult] = useState<GenerateResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [loadingStep, setLoadingStep] = useState(0);

  useEffect(() => {
    if (!loading) return;

    const timer = window.setInterval(() => {
      setLoadingStep((step) => Math.min(step + 1, 4));
    }, 1600);

    return () => window.clearInterval(timer);
  }, [loading]);

  async function handleGenerate() {
    if (!novelText.trim()) {
      setError("请至少输入 3 个章节或分节。");
      return;
    }

    setLoadingStep(1);
    setLoading(true);
    setError("");

    try {
      setResult(await generateScreenplay(novelText));
    } catch (err) {
      const message = err instanceof Error ? err.message : "";
      setError(
        message.includes("at least 3 chapters are required")
          ? "请至少输入 3 个章节或分节。"
          : "生成失败，请检查后端服务、AI 配置或输入章节格式。"
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#18181b",
          colorInfo: "#18181b",
          colorBorder: "#dedede",
          colorText: "#18181b",
          colorTextSecondary: "#666666",
          borderRadius: 8,
          fontFamily: 'Inter, Geist, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif'
        }
      }}
    >
      <main className="page-shell">
        <header className="hero-section">
          <div className="hero-copy">
            <span className="hero-eyebrow">AI STORY ADAPTATION WORKBENCH</span>
            <h1>AI Showrunner Workbench</h1>
            <p className="hero-subtitle">从多章节小说到可校验 YAML 剧本的结构化 AI 工作流</p>
          </div>
          <div className="status-tags" aria-label="系统能力">
            <Tag><i className="status-dot" />Provider Ready</Tag>
            <Tag>Map-Reduce</Tag>
            <Tag>Factual Anchors</Tag>
            <Tag>Fidelity Check</Tag>
          </div>
        </header>

        <WorkflowStepper currentStep={result && !loading ? 5 : loading ? loadingStep : 0} loading={loading} />

        {error ? (
          <Alert
            className="error-card"
            type="error"
            showIcon
            message={error}
            description={
              error === "请至少输入 3 个章节或分节。"
                ? "建议粘贴 3～5 章小说内容，以生成更稳定的 Story Bible 和 YAML 剧本。"
                : "确认服务与输入后，可在左侧重新发起生成。"
            }
            closable
            onClose={() => setError("")}
          />
        ) : null}

        <div className="workspace-grid">
          <aside className="control-column">
            <NovelInput
              value={novelText}
              onChange={(value) => {
                setNovelText(value);
                if (error) setError("");
              }}
              onUseSample={() => {
                setNovelText(sampleText);
                setError("");
              }}
            />
            <GeneratePanel loading={loading} onGenerate={handleGenerate} />
          </aside>

          <section className="result-column">
            <OutputPanel result={result} loading={loading} loadingStep={loadingStep} />
          </section>
        </div>

        {result && !loading ? <ResultSections result={result} detailsOnly /> : null}
      </main>
    </ConfigProvider>
  );
}
