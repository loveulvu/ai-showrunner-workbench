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
	plan, err = editor.PrepareLocalPaths(plan)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for index, clip := range plan.Clips {
		downloaded, err := editor.DownloadClip(ctx, clip)
		if err != nil {
			log.Fatal(err)
		}
		plan.Clips[index] = downloaded
		log.Printf("Downloaded clip: shot_id=%s local_path=%s", downloaded.ShotID, downloaded.LocalPath)
	}

	concatFile, err := editor.BuildConcatList(plan)
	if err != nil {
		log.Fatal(err)
	}
	subtitlesFile, err := editor.BuildSRT(plan)
	if err != nil {
		log.Fatal(err)
	}
	result, err := editor.RunFFmpeg(ctx, plan, concatFile)
	if err != nil {
		log.Fatal(err)
	}
	result.SubtitlesFile = subtitlesFile
	log.Printf("Final demo created: %s", result.OutputFile)
	if result.SubtitlesFile != "" {
		log.Printf("Subtitles created: %s", result.SubtitlesFile)
	}
}
