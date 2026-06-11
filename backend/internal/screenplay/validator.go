package screenplay

import "fmt"

func Validate(input Screenplay) ValidationResult {
	var errors []string

	if input.Title == "" {
		errors = append(errors, "title is required")
	}
	if len(input.SourceChapters) == 0 {
		errors = append(errors, "source_chapters are required")
	}
	if len(input.Characters) == 0 {
		errors = append(errors, "characters are required")
	}
	if len(input.Scenes) == 0 {
		errors = append(errors, "scenes are required")
	}

	for i, chapter := range input.SourceChapters {
		label := fmt.Sprintf("source_chapters[%d]", i)
		if chapter.Number <= 0 {
			errors = append(errors, label+".number must be greater than 0")
		}
		if chapter.Title == "" {
			errors = append(errors, label+".title is required")
		}
	}

	for i, character := range input.Characters {
		label := fmt.Sprintf("characters[%d]", i)
		if character.ID == "" {
			errors = append(errors, label+".id is required")
		}
		if character.Name == "" {
			errors = append(errors, label+".name is required")
		}
		if character.Role == "" {
			errors = append(errors, label+".role is required")
		}
	}

	for i, scene := range input.Scenes {
		label := fmt.Sprintf("scenes[%d]", i)
		if scene.ID == "" {
			errors = append(errors, label+".id is required")
		}
		if scene.Location == "" {
			errors = append(errors, label+".location is required")
		}
		if scene.Time == "" {
			errors = append(errors, label+".time is required")
		}
		if scene.Summary == "" {
			errors = append(errors, label+".summary is required")
		}
		if len(scene.Characters) == 0 {
			errors = append(errors, label+".characters are required")
		}
		if len(scene.Dialogues) == 0 {
			errors = append(errors, label+".dialogues are required")
		}

		for j, dialogue := range scene.Dialogues {
			dialogueLabel := fmt.Sprintf("%s.dialogues[%d]", label, j)
			if dialogue.Character == "" {
				errors = append(errors, dialogueLabel+".character is required")
			}
			if dialogue.Line == "" {
				errors = append(errors, dialogueLabel+".line is required")
			}
		}
	}

	if errors == nil {
		errors = []string{}
	}

	return ValidationResult{
		Passed: len(errors) == 0,
		Errors: errors,
	}
}
