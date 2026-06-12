"use client";

import { useState } from "react";
import { Alert, Button, Card, Tag } from "antd";
import { createVideoTask, getVideoTask, renderEditorDemo } from "@/lib/api";
import type { EditResult, ShowrunnerResult, Shot, VideoResult } from "@/lib/api";

type ShowrunnerStatus = "not-started" | "generating" | "failed" | "ready";
type TaskCreationStatus = "not created" | "creating" | "created" | "failed";

type TaskCreationState = {
  status: TaskCreationStatus;
  error: string;
};

type ShowrunnerOutputProps = {
  result: ShowrunnerResult | null;
  status: ShowrunnerStatus;
  error: string;
  onRetry: () => void;
};

export function ShowrunnerOutput({ result, status, error: showrunnerError, onRetry }: ShowrunnerOutputProps) {
  const [tasks, setTasks] = useState<Record<string, VideoResult>>({});
  const [taskCreation, setTaskCreation] = useState<Record<string, TaskCreationState>>({});
  const [refreshErrors, setRefreshErrors] = useState<Record<string, string>>({});
  const [busyAction, setBusyAction] = useState("");
  const [error, setError] = useState("");
  const [editResult, setEditResult] = useState<EditResult | null>(null);
  const characters = arrayOrEmpty(result?.characters);
  const scenes = arrayOrEmpty(result?.scenes);
  const shots = resolveShots(result);
  const warnings = arrayOrEmpty(result?.warnings);
  const demoShots = shots.slice(0, 3);
  const allCreated = demoShots.length > 0 && demoShots.every((shot) => Boolean(tasks[shot.id]?.task_id));
  const allSucceeded = demoShots.length > 0 && demoShots.every((shot) => tasks[shot.id]?.status === "succeeded" && tasks[shot.id]?.video_url);
  const createdCount = demoShots.filter((shot) => Boolean(tasks[shot.id]?.task_id)).length;
  const failedShots = demoShots.filter((shot) => taskCreation[shot.id]?.status === "failed" && !tasks[shot.id]?.task_id);
  const refreshableShots = demoShots.filter((shot) => Boolean(tasks[shot.id]?.task_id));
  const creationStarted = demoShots.some((shot) => taskCreation[shot.id]?.status && taskCreation[shot.id]?.status !== "not created");

  async function handleCreateDemoTasks() {
    if (!window.confirm("This may consume Wan video credits. Create up to 3 video tasks?")) return;
    await createTasksForShots(demoShots, "create");
  }

  async function handleRetryFailedTasks() {
    if (!window.confirm(`Retry ${failedShots.length} failed video task(s)? This may consume Wan video credits.`)) return;
    await createTasksForShots(failedShots, "retry");
  }

  async function createTasksForShots(shotsToCreate: Shot[], action: "create" | "retry") {
    setBusyAction(action);
    setError("");
    for (const shot of shotsToCreate) {
      if (tasks[shot.id]?.task_id) continue;
      setTaskCreation((current) => ({
        ...current,
        [shot.id]: { status: "creating", error: "" }
      }));
      try {
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
        setTaskCreation((current) => ({
          ...current,
          [shot.id]: { status: "created", error: "" }
        }));
      } catch (err) {
        const message = err instanceof Error ? err.message : "Create video task failed";
        setTaskCreation((current) => ({
          ...current,
          [shot.id]: { status: "failed", error: message }
        }));
      }
    }
    setBusyAction("");
  }

  async function handlePollStatuses() {
    if (refreshableShots.length === 0) return;
    setBusyAction("poll");
    setError("");
    const nextErrors: Record<string, string> = {};
    const refreshedEntries = await Promise.all(
      refreshableShots.map(async (shot) => {
        const task = tasks[shot.id];
        if (task.status === "succeeded") return [shot.id, task] as const;
        try {
          return [shot.id, await getVideoTask(task.task_id)] as const;
        } catch (err) {
          nextErrors[shot.id] = err instanceof Error ? err.message : "Refresh video task status failed";
          return [shot.id, task] as const;
        }
      })
    );
    setTasks((current) => ({
      ...current,
      ...Object.fromEntries(refreshedEntries)
    }));
    setRefreshErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) {
      setError(`${Object.keys(nextErrors).length} video task status refresh failed`);
    }
    setBusyAction("");
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
        <ShowrunnerStatusAlert status={status} error={showrunnerError} onRetry={onRetry} />
        <div className="status-tags showrunner-tags">
          <Tag>Mode: {result?.mode ?? "demo"}</Tag>
          <Tag>Demo shots: 3</Tag>
          <Tag>Full asset generation deferred</Tag>
          <Tag>{characters.length} Characters</Tag>
          <Tag>{scenes.length} Scenes</Tag>
          <Tag>{shots.length} Shots</Tag>
          <Tag>{warnings.length} Warnings</Tag>
        </div>

        {error ? <Alert className="error-card video-task-error" type="error" showIcon message={error} /> : null}

        <div className="video-task-list short-demo-workflow">
          <div className="card-heading">
            <span className="section-kicker">PHASE 6 / SHORT DEMO WORKFLOW</span>
            <h3>Short Demo Workflow</h3>
            {demoShots.length ? (
              <p>Uses the first {demoShots.length} shots. Creating tasks may consume Wan video credits.</p>
            ) : (
              <p>No shots available for short demo</p>
            )}
          </div>
          {demoShots.length ? (
            <>
            <div className="short-demo-actions">
              <Button type="primary" disabled={creationStarted || busyAction !== ""} loading={busyAction === "create"} onClick={handleCreateDemoTasks}>
                Create 3 Video Tasks
              </Button>
              <Button disabled={failedShots.length === 0 || busyAction !== ""} loading={busyAction === "retry"} onClick={handleRetryFailedTasks}>
                Retry Failed Tasks
              </Button>
              <Button disabled={refreshableShots.length === 0 || busyAction !== ""} loading={busyAction === "poll"} onClick={handlePollStatuses}>
                Refresh Status
              </Button>
              <Button disabled={!allSucceeded || busyAction !== ""} loading={busyAction === "render"} onClick={handleRenderDemo}>
                Render Final Demo
              </Button>
            </div>
            <Alert
              type={createdCount === demoShots.length ? "success" : failedShots.length ? "warning" : "info"}
              showIcon
              message={`${createdCount}/${demoShots.length} video tasks created`}
              description={`${refreshableShots.length} task${refreshableShots.length === 1 ? "" : "s"} can be refreshed | ${failedShots.length} failed task${failedShots.length === 1 ? "" : "s"} can be retried`}
            />
            {demoShots.map((shot) => {
              const task = tasks[shot.id];
              const creation = taskCreation[shot.id] ?? { status: "not created", error: "" };
              return (
                <div className="video-task-row" key={shot.id}>
                  <div>
                    <strong>{shot.id}</strong>
                    <small>{summarizePrompt(shot.video_prompt || shot.image_prompt)}</small>
                    <small>Duration: {parseDuration(shot.duration_hint)}s</small>
                    <code>Task: {task?.task_id ?? "not created"}</code>
                    <code>Task status: {task?.status ?? "not available"}</code>
                    <code>Video URL: {safeDisplayURL(task?.video_url)}</code>
                    {creation.error ? <small>Failed: {creation.error}</small> : null}
                    {refreshErrors[shot.id] ? <small>Refresh failed: {refreshErrors[shot.id]}</small> : null}
                  </div>
                  <Tag>{creation.status}</Tag>
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
            </>
          ) : (
            <Alert className="video-task-error" type="info" showIcon message="No shots available for short demo" />
          )}
        </div>

        {result ? <pre className="showrunner-json">{JSON.stringify(result, null, 2)}</pre> : null}
      </Card>
    </section>
  );
}

function ShowrunnerStatusAlert({ status, error, onRetry }: { status: ShowrunnerStatus; error: string; onRetry: () => void }) {
  if (status === "generating") {
    return <Alert type="info" showIcon message="Generating showrunner assets..." />;
  }
  if (status === "failed") {
    return (
      <Alert
        className="error-card showrunner-error"
        type="error"
        showIcon
        message={`Showrunner failed: ${error || "Unknown error"}`}
        action={<Button onClick={onRetry}>Retry Showrunner</Button>}
      />
    );
  }
  if (status === "ready") {
    return <Alert type="success" showIcon message="Showrunner ready" />;
  }
  return <Alert type="info" showIcon message="Not started" />;
}

function resolveShots(result: ShowrunnerResult | null): Shot[] {
  if (Array.isArray(result?.shots) && result.shots.length > 0) {
    return result.shots;
  }

  const shotPrompts: unknown = result?.asset_prompts?.shot_prompts;
  if (!Array.isArray(shotPrompts)) {
    return [];
  }

  return shotPrompts.flatMap((value, index) => {
    if (!value || typeof value !== "object") {
      return [];
    }
    const shot = value as Partial<Shot>;
    return [{
      id: shot.id || `shot_${index + 1}`,
      chapter_number: shot.chapter_number ?? 0,
      scene_id: shot.scene_id ?? "",
      characters: arrayOrEmpty(shot.characters),
      dialogue: arrayOrEmpty(shot.dialogue),
      action: shot.action ?? "",
      camera: shot.camera ?? "",
      background: shot.background ?? "",
      duration_hint: shot.duration_hint ?? "5s",
      image_prompt: shot.image_prompt ?? "",
      video_prompt: shot.video_prompt ?? shot.image_prompt ?? "",
      negative_prompt: shot.negative_prompt ?? defaultNegativePrompt,
      audio_prompt: shot.audio_prompt ?? "",
      character_visuals: arrayOrEmpty(shot.character_visuals),
      scene_visuals: shot.scene_visuals ?? "",
      camera_angle: shot.camera_angle ?? "",
      camera_movement: shot.camera_movement ?? "",
      composition: shot.composition ?? "",
      lighting: shot.lighting ?? "",
      motion: shot.motion ?? "",
      continuity_notes: shot.continuity_notes ?? ""
    }];
  });
}

function arrayOrEmpty<T>(value: T[] | undefined): T[] {
  return Array.isArray(value) ? value : [];
}

function videoPromptForShot(shot: Shot) {
  return {
    shot_id: shot.id,
    model: "",
    prompt: shot.video_prompt || shot.image_prompt,
    negative_prompt: shot.negative_prompt || defaultNegativePrompt,
    duration_seconds: parseDuration(shot.duration_hint),
    aspect_ratio: "16:9",
    subtitle: stringListText(shot.dialogue),
    expected_clip_name: `${shot.id}.mp4`
  };
}

const defaultNegativePrompt = "blurry, distorted face, inconsistent character, different outfit, extra limbs, bad hands, text, subtitles, watermark, logo, low quality, jump cut, flickering, deformed body";

function parseDuration(value: string): number {
  const duration = Number.parseInt(value, 10);
  return Number.isFinite(duration) && duration > 0 ? duration : 5;
}

function summarizePrompt(value: string): string {
  return value.length > 180 ? `${value.slice(0, 180)}...` : value;
}

function safeDisplayURL(value?: string): string {
  if (!value) return "not available";
  try {
    const url = new URL(value);
    return `${url.host}${url.pathname}`;
  } catch {
    return "invalid URL";
  }
}

function stringListText(values: string[]): string {
  return values.join(" ");
}
