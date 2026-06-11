"use client";

import { useState } from "react";
import { Alert, Button, Card, Tag } from "antd";
import { createVideoTask, getVideoTask } from "@/lib/api";
import type { ShowrunnerResult, Shot, VideoResult } from "@/lib/api";

export function ShowrunnerOutput({ result }: { result: ShowrunnerResult }) {
  const [tasks, setTasks] = useState<Record<string, VideoResult>>({});
  const [busyShot, setBusyShot] = useState("");
  const [error, setError] = useState("");
  const eligibleShots = result.shots.filter((shot) => shot.video_prompt || shot.image_prompt);

  async function handleCreate(shot: Shot) {
    setBusyShot(shot.id);
    setError("");
    try {
      const taskID = await createVideoTask({
        shot_id: shot.id,
        model: "mock-video-model",
        prompt: shot.video_prompt || shot.image_prompt,
        negative_prompt: "",
        duration_seconds: parseDuration(shot.duration_hint),
        aspect_ratio: "16:9",
        subtitle: shot.dialogue,
        expected_clip_name: `${shot.id}.mp4`
      });
      setTasks((current) => ({
        ...current,
        [shot.id]: {
          task_id: taskID,
          shot_id: shot.id,
          status: "pending",
          video_url: "",
          error_message: ""
        }
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Create mock video task failed");
    } finally {
      setBusyShot("");
    }
  }

  async function handleRefresh(shot: Shot) {
    const task = tasks[shot.id];
    if (!task) return;
    setBusyShot(shot.id);
    setError("");
    try {
      const refreshed = await getVideoTask(task.task_id);
      setTasks((current) => ({ ...current, [shot.id]: refreshed }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Get mock video task failed");
    } finally {
      setBusyShot("");
    }
  }

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
        {error ? <Alert className="error-card video-task-error" type="error" showIcon message={error} /> : null}
        {eligibleShots.length ? (
          <div className="video-task-list">
            <div className="card-heading">
              <span className="section-kicker">PHASE 3 / MOCK VIDEO TASKS</span>
              <h3>Video Generator Interface</h3>
              <p>创建并查询内存中的 Mock 视频任务，不调用真实视频生成服务。</p>
            </div>
            {eligibleShots.map((shot) => {
              const task = tasks[shot.id];
              return (
                <div className="video-task-row" key={shot.id}>
                  <div>
                    <strong>{shot.id}</strong>
                    <small>{shot.video_prompt || shot.image_prompt}</small>
                    {task ? <code>{task.task_id}</code> : null}
                    {task?.video_url ? <code>{task.video_url}</code> : null}
                  </div>
                  <Tag>{task?.status ?? "not created"}</Tag>
                  {task ? (
                    <Button loading={busyShot === shot.id} onClick={() => handleRefresh(shot)}>Refresh Mock Status</Button>
                  ) : (
                    <Button loading={busyShot === shot.id} onClick={() => handleCreate(shot)}>Create Mock Video Task</Button>
                  )}
                </div>
              );
            })}
          </div>
        ) : null}
        <pre className="showrunner-json">{JSON.stringify(result, null, 2)}</pre>
      </Card>
    </section>
  );
}

function parseDuration(value: string): number {
  const duration = Number.parseInt(value, 10);
  return Number.isFinite(duration) && duration > 0 ? duration : 5;
}
