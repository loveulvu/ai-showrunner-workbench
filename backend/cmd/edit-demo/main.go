package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"ai-showrunner-workbench/internal/editor"
)

func main() {
	planPath := flag.String("plan", "", "path to editing_plan.json")
	flag.Parse()
	if *planPath == "" {
		log.Fatal("--plan is required")
	}

	payload, err := os.ReadFile(*planPath)
	if err != nil {
		log.Fatalf("read editing plan: %v", err)
	}
	var plan editor.EditingPlan
	if err := json.Unmarshal(payload, &plan); err != nil {
		log.Fatalf("parse editing plan: %v", err)
	}
	result, err := editor.Render(context.Background(), plan)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Final demo created: %s", result.OutputFile)
	if result.SubtitlesFile != "" {
		log.Printf("Subtitles created: %s", result.SubtitlesFile)
	}
}
