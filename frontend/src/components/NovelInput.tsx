import { Button, Card, Input } from "antd";

const { TextArea } = Input;

type NovelInputProps = {
  value: string;
  onChange: (value: string) => void;
  onUseSample: () => void;
};

export function NovelInput({ value, onChange, onUseSample }: NovelInputProps) {
  return (
    <Card className="tool-card input-panel">
      <div className="card-heading">
        <div className="heading-line">
          <div>
            <span className="section-kicker">01 / SOURCE NOVEL</span>
            <h2>小说输入</h2>
          </div>
          <span className="input-format">TXT · 3–5 CHAPTERS</span>
        </div>
        <p>输入三个章节以上的小说文本，系统会先进行章节级结构化分析，而不是一次性直接生成剧本。</p>
      </div>
      <TextArea
        aria-label="小说章节输入"
        className="novel-textarea"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder={"第一章 标题\n粘贴章节正文...\n\n第二章 标题\n粘贴章节正文...\n\nChapter 3 Title\n粘贴章节正文..."}
        autoSize={{ minRows: 18, maxRows: 28 }}
      />
      <div className="input-footer">
        <span><strong>{value.length.toLocaleString()}</strong> 字符</span>
        <Button className="secondary-button" type="default" onClick={onUseSample}>
          填入示例小说
        </Button>
      </div>
    </Card>
  );
}
