// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

// Edit Diffsay and run its tests inside the container.
// Init installs Python; Pi then edits const.py and runs pytest.
//
// Prerequisites:
//   - export PROXY_KEY=sk-or-...
//   - clivern/swarm:v0.8.3
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

const diffsayRepo = "https://github.com/Clivern/Diffsay.git"

const initBash = `
apt-get update
apt-get install -y --no-install-recommends python3 python3-pip python3-venv
pip3 install --break-system-packages uv ruff
`

const prompt = `Add "Gemfile.lock" to LOW_PRIORITY_FILES in src/diffsay/const.py if missing and run the tests and provide the full test results`

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
	defer cancel()

	id := uuid.NewString()
	fmt.Printf("job id=%s\ncloning %s …\n", id, diffsayRepo)

	result, err := swarm.Run(ctx, swarm.RunRequest{
		WorkDir:     "/tmp/basement",
		ID:          id,
		RepoURL:     diffsayRepo,
		Prompt:      prompt,
		PIModel:     "openrouter/anthropic/claude-sonnet-4.5",
		ProxyKey:    os.Getenv("PROXY_KEY"),
		ProxyURL:    "http://host.docker.internal:8080/api",
		DockerImage: "clivern/swarm:v0.8.3",
		Container: swarm.Container{
			InitBash: initBash,
			Memory:   "2g",
			CPUs:     "1",
		},
		Cleanup: true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "run failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("=", 72))
	fmt.Println("--- summary ---")
	fmt.Println(strings.TrimSpace(result.Summary))
	fmt.Printf("total_tokens=%d\n", result.TotalTokens)
	fmt.Println("--- changed files ---")
	if len(result.ChangedFiles) == 0 {
		fmt.Println("(none)")
	} else {
		for _, f := range result.ChangedFiles {
			fmt.Printf("%s\t%s\n", f.Status, f.Path)
		}
	}
	fmt.Println("--- patch ---")
	if strings.TrimSpace(result.Patch) == "" {
		fmt.Println("(empty)")
	} else {
		fmt.Print(result.Patch)
		if !strings.HasSuffix(result.Patch, "\n") {
			fmt.Println()
		}
	}
}
