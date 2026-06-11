"use client";

import { useState } from "react";
import { Alert, Button, Card, Tag } from "antd";
import { createVideoTask, getVideoTask, renderEditorDemo } from "@/lib/api";
import type { EditResult, ShowrunnerResult, Shot, VideoResult } from "@/lib/api";

export function ShowrunnerOutput({ result }: { result: ShowrunnerResult }) {
  const [tasks, setTasks] = useState<Record<string, VideoResult>>({});
  const [busyAction, setBusyAction] = useState("");
  const [error, setError] = useState("");
  const [editResult, setEditResult] = useState<EditResult | null>(null);
  const demoShots = result.shots.slice(0, 3);
  const allCreated = demoShots.length > 0 && demoShots.every((shot) => tasks[shot.id]);
  const allSucceeded = demoShots.length > 0 && demoShots.every((shot) => tasks[shot.id]?.status === "succeeded" && tasks[shot.id]?.video_url);

  async function handleCreateDemoTasks() {
    if (!window.confirm("This may consume Wan video credits. Create up to 3 video tasks?")) return;
    setBusyAction("create");
    setError("");
    try {
      for (const shot of demoShots) {
        if (tasks[shot.id]) continue;
        const taskID = await createVideoTask(videoPromptForShot(shot));
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
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Create video tasks failed");
    } finally {
      setBusyAction("");
    }
  }

  async function handlePollStatuses() {
    if (!allCreated) return;
    setBusyAction("poll");
    setError("");
    try {
      let current = { ...tasks };
      while (demoShots.some((shot) => !isTerminal(current[shot.id]?.status))) {
        const refreshedEntries = await Promise.all(
          demoShots.map(async (shot) => {
            const task = current[shot.id];
            if (!task || isTerminal(task.status)) return [shot.id, task] as const;
            return [shot.id, await getVideoTask(task.task_id)] as const;
          })
        );
        current = Object.fromEntries(refreshedEntries) as Record<string, VideoResult>;
        setTasks(current);
        if (demoShots.some((shot) => !isTerminal(current[shot.id]?.status))) {
          await delay(10_000);
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Refresh video task status failed");
    } finally {
      setBusyAction("");
    }
  }

  async function handleRenderDemo() {
    if (!allSucceeded) return;
    setBusyAction("render");
    setError("");
    setEditResult(null);
    try {
      setEditResult(await renderEditorDemo({
        output_file: "../outputs/final_demo.mp4",
        aspect_ratio: "16:9",
        resolution: "1280x720",
        fps: 24,
        clips: demoShots.map((shot) => ({
          shot_id: shot.id,
          source_url: tasks[shot.id].video_url,
          duration_seconds: parseDuration(shot.duration_hint),
          subtitle: stringListText(shot.dialogue)
        }))
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Render final demo failed");
    } finally {
      setBusyAction("");
    }
  }

  return (
    <section className="detail-area">
      <Card className="tool-card artifact-card">
        <div className="card-heading">
          <span className="section-kicker">04 / SHOWRUNNER OUTPUT</span>
          <h2>Showrunner Output</h2>
          <p>Characters, scenes, chapter breakdowns, shots, and generation prompts.</p>
        </div>
        <div className="status-tags showrunner-tags">
          <Tag>{result.characters.length} Characters</Tag>
          <Tag>{result.scenes.length} Scenes</Tag>
          <Tag>{result.shots.length} Shots</Tag>
          <Tag>{result.warnings.length} Warnings</Tag>
        </div>

        {error ? <Alert className="error-card video-task-error" type="error" showIcon message={error} /> : null}

        {result.shots.length ? (
          <div className="video-task-list short-demo-workflow">
            <div className="card-heading">
              <span className="section-kicker">PHASE 6 / SHORT DEMO WORKFLOW</span>
              <h3>Generate Short Demo</h3>
              <p>Uses the first {demoShots.length} shots. Creating tasks may consume Wan video credits.</p>
            </div>
            <div className="short-demo-actions">
              <Button type="primary" disabled={allCreated || busyAction !== ""} loading={busyAction === "create"} onClick={handleCreateDemoTasks}>
                Create 3 Video Tasks
              </Button>
              <Button disabled={!allCreated || busyAction !== ""} loading={busyAction === "poll"} onClick={handlePollStatuses}>
                Refresh Status
              </Button>
              <Button disabled={!allSucceeded || busyAction !== ""} loading={busyAction === "render"} onClick={handleRenderDemo}>
                Render Final Demo
              </Button>
            </div>
            {demoShots.map((shot) => {
              const task = tasks[shot.id];
              return (
                <div className="video-task-row" key={shot.id}>
                  <div>
                    <strong>{shot.id}</strong>
                    <small>{summarizePrompt(shot.video_prompt || shot.image_prompt)}</small>
                    <small>Duration: {parseDuration(shot.duration_hint)}s</small>
                    <code>Task: {task?.task_id ?? "not created"}</code>
                    <code>Video URL: {task?.video_url || "not available"}</code>
                  </div>
                  <Tag>{task?.status ?? "not created"}</Tag>
                </div>
              );
            })}
            {editResult ? (
              <Alert
                type="success"
                showIcon
                message="Final demo rendered"
                description={`Output: ${editResult.output_file}${editResult.subtitles_file ? ` | Subtitles: ${editResult.subtitles_file}` : ""}`}
              />
            ) : null}
          </div>
        ) : (
          <Alert className="video-task-error" type="info" showIcon message="No shots available for short demo." />
        )}

        <pre className="showrunner-json">{JSON.stringify(result, null, 2)}</pre>
      </Card>
    </section>
  );
}

function videoPromptForShot(shot: Shot) {
  return {
    shot_id: shot.id,
    model: "",
    prompt: shot.video_prompt || shot.image_prompt,
    negative_prompt: "",
    duration_seconds: parseDuration(shot.duration_hint),
    aspect_ratio: "16:9",
    subtitle: stringListText(shot.dialogue),
    expected_clip_name: `${shot.id}.mp4`
  };
}

function parseDuration(value: string): number {
  const duration = Number.parseInt(value, 10);
  return Number.isFinite(duration) && duration > 0 ? duration : 5;
}

function summarizePrompt(value: string): string {
  return value.length > 180 ? `${value.slice(0, 180)}...` : value;
}

function isTerminal(status?: string): boolean {
  return status === "succeeded" || status === "failed";
}

function delay(milliseconds: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, milliseconds));
}

function stringListText(values: string[]): string {
  return values.join(" ");
}
