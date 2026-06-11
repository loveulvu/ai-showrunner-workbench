"use client";

import { Card, Tag } from "antd";
import type { ShowrunnerResult } from "@/lib/api";

export function ShowrunnerOutput({ result }: { result: ShowrunnerResult }) {
  return (
    <section className="detail-area">
      <Card className="tool-card artifact-card">
        <div className="card-heading">
          <span className="section-kicker">04 / SHOWRUNNER OUTPUT</span>
          <h2>Showrunner 输出</h2>
          <p>角色设定、场景设定、章节拆解、分镜列表与资产生成提示词。</p>
        </div>
        <div className="status-tags showrunner-tags">
          <Tag>{result.characters.length} Characters</Tag>
          <Tag>{result.scenes.length} Scenes</Tag>
          <Tag>{result.shots.length} Shots</Tag>
          <Tag>{result.warnings.length} Warnings</Tag>
        </div>
        <pre className="showrunner-json">{JSON.stringify(result, null, 2)}</pre>
      </Card>
    </section>
  );
}
