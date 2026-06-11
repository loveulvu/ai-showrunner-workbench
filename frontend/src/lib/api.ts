export type CharacterMention = {
  name: string;
  role_in_chapter: string;
  traits: string[];
  state_change: string;
};

export type SceneCandidate = {
  location: string;
  time: string;
  purpose: string;
  characters: string[];
  key_events: string[];
};

export type ChapterAnalysis = {
  chapter_number: number;
  chapter_title: string;
  summary: string;
  characters: CharacterMention[];
  locations: string[];
  key_events: string[];
  conflicts: string[];
  scene_candidates: SceneCandidate[];
  factual_anchors: string[];
};

export type StoryCharacter = {
  id: string;
  name: string;
  role: string;
  motivation: string;
};

export type TimelineEvent = {
  chapter_number: number;
  event: string;
};

export type ScenePlanItem = {
  id: string;
  source_chapter: number;
  summary: string;
  location: string;
  time: string;
  characters: string[];
};

export type StoryBible = {
  title: string;
  logline: string;
  global_characters: StoryCharacter[];
  timeline: TimelineEvent[];
  main_conflict: string;
  scene_plan: ScenePlanItem[];
};

export type SourceChapter = {
  number: number;
  title: string;
  summary: string;
};

export type ScreenplayCharacter = {
  id: string;
  name: string;
  role: string;
  description: string;
};

export type Dialogue = {
  character: string;
  emotion: string;
  line: string;
};

export type Scene = {
  id: string;
  source_chapter: number;
  location: string;
  time: string;
  summary: string;
  characters: string[];
  dialogues: Dialogue[];
  actions: string[];
};

export type Screenplay = {
  title: string;
  source_chapters: SourceChapter[];
  characters: ScreenplayCharacter[];
  scenes: Scene[];
};

export type Validation = {
  passed: boolean;
  errors: string[];
};

export type FidelityIssue = {
  field: string;
  severity: "low" | "medium" | "high" | string;
  problem: string;
  suggestion: string;
};

export type FidelityResult = {
  passed: boolean;
  issues: FidelityIssue[];
};

export type GenerateMeta = {
  ai_provider: "mock" | "real" | string;
  ai_model: string;
};

export type GenerateResponse = {
  chapter_count: number;
  chapter_analyses: ChapterAnalysis[];
  story_bible: StoryBible;
  screenplay_json: Screenplay;
  screenplay_yaml: string;
  validation: Validation;
  fidelity_result: FidelityResult;
  meta: GenerateMeta;
};

export type CharacterProfile = {
  id: string;
  name: string;
  role: string;
  personality: string[];
  appearance: string;
  costume: string;
  voice_style: string;
  key_motivation: string;
  consistency_notes: string[];
};

export type SceneProfile = {
  id: string;
  name: string;
  location: string;
  time_of_day: string;
  atmosphere: string;
  visual_style: string;
  key_props: string[];
  consistency_notes: string[];
};

export type ChapterBreakdown = {
  chapter_number: number;
  chapter_title: string;
  summary: string;
  main_characters: string[];
  main_scenes: string[];
  emotional_arc: string;
  key_events: string[];
};

export type Shot = {
  id: string;
  chapter_number: number;
  scene_id: string;
  characters: string[];
  dialogue: string;
  action: string;
  camera: string;
  background: string;
  duration_hint: string;
  image_prompt: string;
  video_prompt: string;
  audio_prompt: string;
};

export type AssetPromptSet = {
  character_prompts: Record<string, string>;
  background_prompts: Record<string, string>;
  shot_prompts: Record<string, string>;
  voice_prompts: Record<string, string>;
};

export type ShowrunnerResult = {
  characters: CharacterProfile[];
  scenes: SceneProfile[];
  chapters: ChapterBreakdown[];
  shots: Shot[];
  asset_prompts: AssetPromptSet;
  warnings: string[];
};

export type VideoTaskStatus = "pending" | "running" | "succeeded" | "failed";

export type VideoPrompt = {
  shot_id: string;
  model: string;
  prompt: string;
  negative_prompt: string;
  duration_seconds: number;
  aspect_ratio: string;
  subtitle: string;
  expected_clip_name: string;
};

export type VideoResult = {
  task_id: string;
  shot_id: string;
  status: VideoTaskStatus;
  video_url: string;
  error_message: string;
};

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function generateScreenplay(novelText: string): Promise<GenerateResponse> {
  let response: Response;

  try {
    response = await fetch(`${API_BASE_URL}/api/generate`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        novel_text: novelText
      })
    });
  } catch {
    throw new Error("请求失败，请确认后端服务已启动");
  }

  const data = await response.json().catch(() => null);
  if (!response.ok) {
    throw new Error(data?.error ?? "生成失败，请稍后重试");
  }

  return data as GenerateResponse;
}

export async function generateShowrunner(
  result: GenerateResponse,
  style = "",
  language = "zh-CN"
): Promise<ShowrunnerResult> {
  let response: Response;

  try {
    response = await fetch(`${API_BASE_URL}/api/showrunner/generate`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        screenplay: result.screenplay_json,
        story_bible: result.story_bible,
        chapters: result.chapter_analyses,
        style,
        language
      })
    });
  } catch {
    throw new Error("Showrunner request failed");
  }

  const data = await response.json().catch(() => null);
  if (!response.ok) {
    throw new Error(data?.error ?? "Showrunner generation failed");
  }

  return data as ShowrunnerResult;
}

export async function createVideoTask(prompt: VideoPrompt): Promise<string> {
  const response = await fetch(`${API_BASE_URL}/api/video/tasks`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(prompt)
  });
  const data = await response.json().catch(() => null);
  if (!response.ok) {
    throw new Error(data?.error ?? "Create mock video task failed");
  }
  return data.task_id as string;
}

export async function getVideoTask(taskID: string): Promise<VideoResult> {
  const response = await fetch(`${API_BASE_URL}/api/video/tasks/${encodeURIComponent(taskID)}`);
  const data = await response.json().catch(() => null);
  if (!response.ok) {
    throw new Error(data?.error ?? "Get mock video task failed");
  }
  return data as VideoResult;
}
