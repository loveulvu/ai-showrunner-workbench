package showrunner

import (
	"encoding/json"
	"fmt"
)

func BuildPrompt(input GenerateInput) string {
	if input.Style == "" {
		input.Style = "cinematic xianxia short drama"
	}
	payload, _ := json.MarshalIndent(input, "", "  ")
	return fmt.Sprintf(`You are a storyboard director and video prompt engineer for cinematic short drama.
Generate concise, structured production-planning assets optimized for Wan video generation.

Return one strict JSON object with exactly these top-level fields:
- characters
- scenes
- chapters
- shots
- asset_prompts
- warnings

Requirements:
1. Build a continuity bible. Every main character visual_identity must define age, face, hairstyle, costume, color_palette, expression_baseline, body_type, and a concise consistency_prompt.
2. Every scene visual_identity must define architecture, lighting, color_palette, atmosphere, key_props, and a concise consistency_prompt.
3. Every shot must contain action or dialogue.
4. Every shot must contain a non-empty image_prompt.
5. Every shot must reference the relevant stable character and scene identity through character_visuals and scene_visuals.
6. Generate production prompts only. Do not claim to generate real images, video, audio, or edited footage.
7. Preserve the facts, character identities, chapter ownership, and events in the input.
8. Output JSON only. Do not output Markdown, YAML, explanations, or comments.
9. Use the requested style and language when they are provided.
10. Every list field must be a JSON array, even when it contains only one item or no items.
11. Never output an array field as a string.
12. The output must be strict JSON only, with no Markdown fences and no explanatory text.
13. The first three shots must form one continuous mini-scene in the same location and time: setup/entrance, conflict or key action, then reaction/result. Preserve positions, costume, lighting, props, and screen direction between them.
14. Each shot must define camera_angle, camera_movement, composition, lighting, motion, and continuity_notes. Use controlled camera movement and realistic temporal motion, not a static illustration.
15. Each video_prompt must be concise and include: style, fixed character face/hair/costume, fixed scene, specific action, camera movement, lighting, composition, stable face, stable costume, realistic motion, no text, no subtitles, no watermark.
16. Default style is "cinematic xianxia short drama". Do not use ink painting unless the requested style explicitly asks for it.
17. Each negative_prompt must include: blurry, distorted face, inconsistent character, different outfit, extra limbs, bad hands, text, subtitles, watermark, logo, low quality, jump cut, flickering, deformed body.
18. Keep each video_prompt focused and reasonably short; avoid repetitive prose.

List fields include:
personality, appearance, costume, voice_style, key_motivation, consistency_notes, key_props,
main_characters, main_scenes, key_events, characters, dialogue, character_visuals, and warnings.

Required character fields:
id, name, role, personality, appearance, costume, voice_style, key_motivation, consistency_notes, visual_identity

Required scene fields:
id, name, location, time_of_day, atmosphere, visual_style, key_props, consistency_notes, visual_identity

Required chapter fields:
chapter_number, chapter_title, summary, main_characters, main_scenes, emotional_arc, key_events

Required shot fields:
id, chapter_number, scene_id, characters, dialogue, action, camera, background, duration_hint,
character_visuals, scene_visuals, camera_angle, camera_movement, composition, lighting, motion,
continuity_notes, image_prompt, video_prompt, negative_prompt, audio_prompt

Required asset_prompts fields:
character_prompts, background_prompts, shot_prompts, voice_prompts

Input JSON:
%s`, string(payload))
}

func BuildRepairPrompt(raw string, reason string) string {
	return fmt.Sprintf(`Repair the following animation showrunner output.
Return one strict JSON object only, with no Markdown fences, explanations, YAML, or comments.
Preserve all usable data. Ensure the top-level "shots" field is a non-empty JSON array.
Preserve or restore character and scene visual_identity, shot continuity fields, concise video_prompt, and negative_prompt.
Normalize identifiers and chapter numbers to strings or numbers accepted by JSON.
The previous output failed because: %s

Previous output:
%s`, reason, raw)
}
