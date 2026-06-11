"use client";

import type { ReactNode } from "react";
import { Button, Card, message } from "antd";

type YamlPreviewProps = {
  yaml: string;
  embedded?: boolean;
};

export function YamlPreview({ yaml, embedded = false }: YamlPreviewProps) {
  const [messageApi, contextHolder] = message.useMessage();
  const lines = (yaml || "# 生成后将在这里展示结构化剧本 YAML").split("\n");

  async function handleCopy() {
    if (!yaml) return;
    await navigator.clipboard.writeText(yaml);
    messageApi.success("YAML 已复制");
  }

  function handleDownload() {
    if (!yaml) return;
    const blob = new Blob([yaml], { type: "text/yaml;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "screenplay.yaml";
    link.click();
    URL.revokeObjectURL(url);
  }

  const content = (
    <>
      {contextHolder}
      <div className="preview-actions">
        <div>
          <span className="section-kicker">YAML EXPORT</span>
          <h2>结构化剧本</h2>
          <p className="module-explainer">最终输出为 YAML，便于复制、下载、编辑和继续二次创作。</p>
        </div>
        {yaml ? (
          <div className="button-row">
            <Button className="secondary-button" onClick={handleCopy}>复制 YAML</Button>
            <Button className="secondary-button" onClick={handleDownload}>下载 screenplay.yaml</Button>
          </div>
        ) : null}
      </div>
      <div className="code-window">
        <div className="code-toolbar">
          <div className="window-dots" aria-hidden="true"><i /><i /><i /></div>
          <span>screenplay.yaml</span>
          <small>{lines.length} lines · UTF-8</small>
        </div>
        <pre className="yaml-output"><code>{lines.map((line, index) => (
          <span className="code-line" key={`${index}-${line}`}>
            <span className="line-number">{String(index + 1).padStart(2, "0")}</span>
            <span className="line-content">{highlightYamlLine(line)}</span>
          </span>
        ))}</code></pre>
      </div>
    </>
  );

  return embedded ? <div className="yaml-panel">{content}</div> : <Card className="tool-card yaml-panel">{content}</Card>;
}

function highlightYamlLine(line: string): ReactNode {
  if (line.trimStart().startsWith("#")) return <span className="yaml-comment">{line}</span>;

  const match = line.match(/^(\s*(?:-\s+)?)([\w-]+)(:)(.*)$/);
  if (!match) return line;

  return (
    <>
      {match[1]}<span className="yaml-key">{match[2]}</span><span className="yaml-punctuation">{match[3]}</span><span className="yaml-value">{match[4]}</span>
    </>
  );
}
