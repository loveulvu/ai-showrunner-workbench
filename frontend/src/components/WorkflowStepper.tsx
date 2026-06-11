type WorkflowStepperProps = {
  currentStep: number;
  loading: boolean;
};

const steps = [
  { label: "小说输入", meta: "Source" },
  { label: "章节分析", meta: "Map" },
  { label: "Story Bible", meta: "Reduce" },
  { label: "事实检查", meta: "Verify" },
  { label: "YAML 导出", meta: "Export" }
];

export function WorkflowStepper({ currentStep, loading }: WorkflowStepperProps) {
  return (
    <nav className="workflow-stepper" aria-label="生成流程">
      {steps.map((step, index) => {
        const state = currentStep > index ? "complete" : currentStep === index ? "active" : "pending";

        return (
          <div className={`workflow-step workflow-step-${state}`} key={step.label} aria-current={state === "active" ? "step" : undefined}>
            <div className="step-marker" aria-hidden="true">
              {state === "complete" ? (
                <svg viewBox="0 0 20 20">
                  <path d="m5 10.5 3 3 7-7" />
                </svg>
              ) : (
                <span>{String(index + 1).padStart(2, "0")}</span>
              )}
            </div>
            <div className="step-copy">
              <strong>{step.label}</strong>
              <span>{state === "active" && loading ? "Processing" : step.meta}</span>
            </div>
            {index < steps.length - 1 ? <i className="step-connector" aria-hidden="true" /> : null}
          </div>
        );
      })}
    </nav>
  );
}
