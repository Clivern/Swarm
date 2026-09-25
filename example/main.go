// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

// Five concurrent agent runs against github.com/Clivern/Diffsay.
//
// Prerequisites:
//   - docker build -t swarm-pi:local ..   (from repo root)
//   - export OPENROUTER_API_KEY=sk-or-...
//
// Run:
//
//	go run ./
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/clivern/swarm"
	"github.com/google/uuid"
)

var tasks = []string{
	`Add "Gemfile.lock" to LOW_PRIORITY_FILES in src/diffsay/const.py if missing.
Only edit that file.`,

	`Add "Poetry.lock" to LOW_PRIORITY_FILES in src/diffsay/const.py if missing.
Only edit that file.`,

	`In tests/test_cli.py add function test_uv_lock_in_low_priority_files:
import LOW_PRIORITY_FILES from diffsay.const and assert "uv.lock" in LOW_PRIORITY_FILES.
Only edit tests/test_cli.py.`,

	`Add a one-sentence docstring to the module at the top of src/diffsay/const.py
explaining LOW_PRIORITY_FILES.`,

	`In README.md under ### Develop add this line immediately after the heading (before the code block):
"Use uv to sync dependencies and run checks."
Do not skip; README must change even if pytest is mentioned in the code block below.`,
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var failed int

	for i, prompt := range tasks {
		wg.Add(1)
		go func(n int, prompt string) {
			defer wg.Done()

			id := uuid.NewString()
			log.Printf("[task %d] start id=%s", n+1, id)

			result, err := swarm.Run(ctx, swarm.RunRequest{
				WorkDir:          ".work",
				ID:               id,
				RepoURL:          "https://github.com/Clivern/Diffsay.git",
				Prompt:           prompt,
				PIModel:          "openrouter/anthropic/claude-sonnet-4.5",
				OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
				DockerImage:      "swarm:v0.1.0",
			})
			if err != nil {
				mu.Lock()
				failed++
				mu.Unlock()
				log.Printf("[task %d] error id=%s: %v", n+1, id, err)
				return
			}

			mu.Lock()
			PrintTaskResult(n+1, id, result)
			mu.Unlock()
		}(i, prompt)
	}

	wg.Wait()
	if failed > 0 {
		log.Fatalf("%d of %d tasks failed", failed, len(tasks))
	}
	fmt.Println("all tasks finished")
}

func PrintTaskResult(n int, id string, result *swarm.Result) {
	bar := strings.Repeat("=", 72)
	fmt.Printf("\n%s\n[task %d] id=%s\nout: %s\n%s\n", bar, n, id, result.OutDir, bar)
	fmt.Println("--- summary ---")
	fmt.Println(strings.TrimSpace(result.Summary))
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
