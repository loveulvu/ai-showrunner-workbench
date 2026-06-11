import { Button, Spin } from "antd";

type GeneratePanelProps = {
  loading: boolean;
  onGenerate: () => void;
};

export function GeneratePanel({ loading, onGenerate }: GeneratePanelProps) {
  return (
    <div className="generate-panel">
      <Button className="primary-button" type="primary" onClick={onGenerate} disabled={loading}>
        <span>{loading ? "正在执行 AI Workflow" : "生成剧本 YAML"}</span>
        <svg viewBox="0 0 20 20" aria-hidden="true">
          <path d="M4 10h12m-4-4 4 4-4 4" />
        </svg>
      </Button>
      {loading ? (
        <div className="loading-row">
          <Spin size="small" />
          <p className="loading-text">正在分析章节、合并故事资料并执行质量检查，请稍候。</p>
        </div>
      ) : (
        <p className="helper-text">生成过程将调用后端配置的 AI provider。</p>
      )}
    </div>
  );
}
