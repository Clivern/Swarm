// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

// Five concurrent agent runs against github.com/Clivern/Diffsay.
//
// Prerequisites:
//   - docker build -t swarm:v0.5.0 ..   (from repo root)
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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
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

type TaskOutcome struct {
	n      int
	id     string
	result *swarm.Result
	err    error
}

type Progress struct {
	mu       sync.Mutex
	inFlight int
	ok       int
	failed   int
	spinner  *spinner.Spinner
}

func (p *Progress) StartTask() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.inFlight++
	p.RefreshLocked()
}

func (p *Progress) FinishTask(ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.inFlight--
	if ok {
		p.ok++
	} else {
		p.failed++
	}
	p.RefreshLocked()
}

func (p *Progress) RefreshLocked() {
	if p.spinner == nil {
		return
	}
	pending := len(tasks) - p.inFlight - p.ok - p.failed
	p.spinner.Suffix = fmt.Sprintf(
		" %d running · %d done · %d failed · %d pending",
		p.inFlight, p.ok, p.failed, pending,
	)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " starting..."
	s.Start()

	prog := &Progress{spinner: s}
	var wg sync.WaitGroup
	var mu sync.Mutex
	outcomes := make([]TaskOutcome, 0, len(tasks))

	for i, prompt := range tasks {
		wg.Add(1)
		go func(n int, prompt string) {
			defer wg.Done()

			id := uuid.NewString()
			prog.StartTask()

			result, err := swarm.Run(ctx, swarm.RunRequest{
				WorkDir:          ".work",
				ID:               id,
				RepoURL:          "https://github.com/Clivern/Diffsay.git",
				Prompt:           prompt,
				PIModel:          "openrouter/anthropic/claude-sonnet-4.5",
				OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
				DockerImage:      "swarm:v0.5.0",
			})

			mu.Lock()
			outcomes = append(outcomes, TaskOutcome{n: n + 1, id: id, result: result, err: err})
			mu.Unlock()

			prog.FinishTask(err == nil)
		}(i, prompt)
	}

	wg.Wait()
	s.Stop()
	fmt.Println()

	sort.Slice(outcomes, func(i, j int) bool { return outcomes[i].n < outcomes[j].n })

	var failed int
	for _, o := range outcomes {
		if o.err != nil {
			failed++
			fmt.Printf("[task %d] id=%s error: %v\n", o.n, o.id, o.err)
			continue
		}
		PrintTaskResult(o.n, o.id, o.result)
	}

	if failed > 0 {
		fmt.Printf("\n%d of %d tasks failed\n", failed, len(tasks))
		os.Exit(1)
	}
	fmt.Println("all tasks finished")
}

func PrintTaskResult(n int, id string, result *swarm.Result) {
	bar := strings.Repeat("=", 72)
	fmt.Printf("\n%s\n[task %d] id=%s\nout: %s\n%s\n", bar, n, id, result.OutDir, bar)
	fmt.Println("--- summary ---")
	fmt.Println(strings.TrimSpace(result.Summary))
	fmt.Println("--- changed files ---")
	if len(result.ChangedFiles) == 0 {
		fmt.Println("(none)")
	} else {
		for _, f := range result.ChangedFiles {
			fmt.Printf("%s\t%s (%d bytes)\n", f.Status, f.Path, len(f.Content))
			if strings.TrimSpace(f.Content) == "" {
				fmt.Println("(no content)")
			} else {
				fmt.Println(f.Content)
			}
			fmt.Println()
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
