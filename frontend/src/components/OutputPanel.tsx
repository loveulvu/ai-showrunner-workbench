"use client";

import { useEffect, useState } from "react";
import { Card, Segmented } from "antd";
import { OverviewContent } from "@/components/ResultSections";
import { ValidationResult } from "@/components/ValidationResult";
import { YamlPreview } from "@/components/YamlPreview";
import type { GenerateResponse } from "@/lib/api";

type OutputSection = "Overview" | "Quality" | "YAML";

type OutputPanelProps = {
  result: GenerateResponse | null;
  loading: boolean;
  loadingStep: number;
};

const loadingStages = [
  "读取章节边界与输入结构",
  "提取人物、事件与事实锚点",
  "合并全局 Story Bible",
  "执行 Fidelity 与 Schema 检查",
  "导出结构化 YAML 剧本"
];

export function OutputPanel({ result, loading, loadingStep }: OutputPanelProps) {
  const [section, setSection] = useState<OutputSection>("Overview");

  useEffect(() => {
    if (result) setSection("Overview");
  }, [result]);

  return (
    <Card className={`tool-card output-card output-card-${loading ? "loading" : result ? section.toLowerCase() : "empty"}`}>
      <div className="output-card-header">
        <div>
          <span className="section-kicker">02 / RESULT CONSOLE</span>
          <h2>生成结果</h2>
        </div>
        <Segmented<OutputSection>
          options={["Overview", "Quality", "YAML"]}
          value={section}
          onChange={setSection}
          disabled={!result || loading}
        />
      </div>

      <div className="output-card-content">
        {loading ? (
          <LoadingState activeStep={loadingStep} />
        ) : !result ? (
          <EmptyState />
        ) : section === "Overview" ? (
          <OverviewContent result={result} />
        ) : section === "Quality" ? (
          <ValidationResult validation={result.validation} fidelityResult={result.fidelity_result} embedded />
        ) : (
          <YamlPreview yaml={result.screenplay_yaml} embedded />
        )}
      </div>
    </Card>
  );
}

function EmptyState() {
  return (
    <div className="result-empty">
      <div className="empty-visual" aria-hidden="true">
        <span />
        <span />
        <span />
        <i />
      </div>
      <span className="empty-index">WORKSPACE READY</span>
      <h3>等待生成结果</h3>
      <p>填入小说并启动工作流后，这里会成为结构化结果控制台。</p>
      <div className="empty-flow">
        {["章节分析", "事实锚点", "Story Bible", "质量检查", "YAML 剧本"].map((item, index) => (
          <span key={item}><b>{String(index + 1).padStart(2, "0")}</b>{item}</span>
        ))}
      </div>
    </div>
  );
}

function LoadingState({ activeStep }: { activeStep: number }) {
  return (
    <div className="result-loading">
      <div className="loading-orbit" aria-hidden="true"><i /></div>
      <span className="section-kicker">PIPELINE RUNNING</span>
      <h3>正在构建结构化剧本</h3>
      <p>结果将逐阶段汇总，无需停留在当前页面。</p>
      <div className="loading-stage-list">
        {loadingStages.map((stage, index) => {
          const state = activeStep > index ? "complete" : activeStep === index ? "active" : "pending";
          return (
            <div className={`loading-stage loading-stage-${state}`} key={stage}>
              <span>{state === "complete" ? "✓" : String(index + 1).padStart(2, "0")}</span>
              <strong>{stage}</strong>
              <i />
            </div>
          );
        })}
      </div>
    </div>
  );
}
