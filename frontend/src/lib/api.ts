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
