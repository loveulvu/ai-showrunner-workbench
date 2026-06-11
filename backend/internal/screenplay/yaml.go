package screenplay

import (
	"strconv"
	"strings"
)

func ToYAML(input Screenplay) string {
	var b strings.Builder

	writeLine(&b, 0, "title: "+quote(input.Title))
	writeLine(&b, 0, "source_chapters:")
	for _, chapter := range input.SourceChapters {
		writeLine(&b, 2, "- number: "+strconv.Itoa(chapter.Number))
		writeLine(&b, 4, "title: "+quote(chapter.Title))
		writeLine(&b, 4, "summary: "+quote(chapter.Summary))
	}

	writeLine(&b, 0, "characters:")
	for _, character := range input.Characters {
		writeLine(&b, 2, "- id: "+quote(character.ID))
		writeLine(&b, 4, "name: "+quote(character.Name))
		writeLine(&b, 4, "role: "+quote(character.Role))
		writeLine(&b, 4, "description: "+quote(character.Description))
	}

	writeLine(&b, 0, "scenes:")
	for _, scene := range input.Scenes {
		writeLine(&b, 2, "- id: "+quote(scene.ID))
		writeLine(&b, 4, "source_chapter: "+strconv.Itoa(scene.SourceChapter))
		writeLine(&b, 4, "location: "+quote(scene.Location))
		writeLine(&b, 4, "time: "+quote(scene.Time))
		writeLine(&b, 4, "summary: "+quote(scene.Summary))

		writeLine(&b, 4, "characters:")
		for _, character := range scene.Characters {
			writeLine(&b, 6, "- "+quote(character))
		}

		writeLine(&b, 4, "dialogues:")
		for _, dialogue := range scene.Dialogues {
			writeLine(&b, 6, "- character: "+quote(dialogue.Character))
			writeLine(&b, 8, "emotion: "+quote(dialogue.Emotion))
			writeLine(&b, 8, "line: "+quote(dialogue.Line))
		}

		writeLine(&b, 4, "actions:")
		for _, action := range scene.Actions {
			writeLine(&b, 6, "- "+quote(action))
		}
	}

	return b.String()
}

func writeLine(b *strings.Builder, indent int, text string) {
	b.WriteString(strings.Repeat(" ", indent))
	b.WriteString(text)
	b.WriteByte('\n')
}

func quote(value string) string {
	return strconv.Quote(value)
}
