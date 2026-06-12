package showrunner

import (
	"encoding/json"
	"fmt"
)

func BuildPrompt(input GenerateInput) string {
	payload, _ := json.MarshalIndent(input, "", "  ")
	return fmt.Sprintf(`You are an animation pre-production showrunner assistant.
Generate structured production-planning assets from the provided screenplay, story bible, and chapter analyses.

Return one strict JSON object with exactly these top-level fields:
- characters
- scenes
- chapters
- shots
- asset_prompts
- warnings

Requirements:
1. Character profiles must remain visually and behaviorally consistent across all shots.
2. Scene profiles must be reusable and use stable scene ids.
3. Every shot must contain action or dialogue.
4. Every shot must contain a non-empty image_prompt.
5. video_prompt and audio_prompt may be empty, but add a clear warning for each missing prompt.
6. Generate production prompts only. Do not claim to generate real images, video, audio, or edited footage.
7. Preserve the facts, character identities, chapter ownership, and events in the input.
8. Output JSON only. Do not output Markdown, YAML, explanations, or comments.
9. Use the requested style and language when they are provided.
10. Every list field must be a JSON array, even when it contains only one item or no items.
11. Never output an array field as a string.
12. The output must be strict JSON only, with no Markdown fences and no explanatory text.

List fields include:
personality, appearance, costume, voice_style, key_motivation, consistency_notes, key_props,
main_characters, main_scenes, key_events, characters, dialogue, and warnings.

Required character fields:
id, name, role, personality, appearance, costume, voice_style, key_motivation, consistency_notes

Required scene fields:
id, name, location, time_of_day, atmosphere, visual_style, key_props, consistency_notes

Required chapter fields:
chapter_number, chapter_title, summary, main_characters, main_scenes, emotional_arc, key_events

Required shot fields:
id, chapter_number, scene_id, characters, dialogue, action, camera, background, duration_hint, image_prompt, video_prompt, audio_prompt

Required asset_prompts fields:
character_prompts, background_prompts, shot_prompts, voice_prompts

Input JSON:
%s`, string(payload))
}

func BuildRepairPrompt(raw string, reason string) string {
	return fmt.Sprintf(`Repair the following animation showrunner output.
Return one strict JSON object only, with no Markdown fences, explanations, YAML, or comments.
Preserve all usable data. Ensure the top-level "shots" field is a non-empty JSON array.
Normalize identifiers and chapter numbers to strings or numbers accepted by JSON.
The previous output failed because: %s

Previous output:
%s`, reason, raw)
}
