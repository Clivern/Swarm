// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

// Read-only codebase question against github.com/Clivern/Ziee (PR triage audit).
//
// Prerequisites:
//   - docker build -t swarm:v0.8.1 ..   (from repo root)
//   - export OPENROUTER_API_KEY=sk-or-...
//
// Run:
//
//	go run ./
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/clivern/swarm"
	"github.com/google/uuid"
)

const zieeRepo = "https://github.com/Clivern/Ziee.git"

const prompt = "Does the PR triage functionality is done?"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
	defer cancel()

	id := uuid.NewString()
	fmt.Printf("job id=%s\ncloning %s …\n", id, zieeRepo)

	result, err := swarm.Run(ctx, swarm.RunRequest{
		WorkDir:          "/tmp/basement",
		ID:               id,
		RepoURL:          zieeRepo,
		Prompt:           prompt,
		PIModel:          "openrouter/anthropic/claude-sonnet-4.5",
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		DockerImage:      "swarm:v0.8.1",
		Cleanup:          true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "run failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("=", 72))
	fmt.Println("--- answer (result.Summary) ---")
	fmt.Println(strings.TrimSpace(result.Summary))
	fmt.Printf("total_tokens=%d\n", result.TotalTokens)

	if strings.TrimSpace(result.Patch) != "" || len(result.ChangedFiles) > 0 {
		fmt.Println("\nwarning: agent modified the tree")
		fmt.Printf("changed_files=%d patch_bytes=%d\n", len(result.ChangedFiles), len(result.Patch))
	} else {
		fmt.Println("\n(no repo changes)")
	}
}
