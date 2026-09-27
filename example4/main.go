// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

// Context deadline stops a long init script.
// Init sleeps for 10 minutes. The deadline is 90 seconds and includes the clone.
// A passing run prints "init sleep start", then returns a docker error before the sleep finishes.
//
// Prerequisites:
//   - clivern/swarm:v0.8.5 already pulled
//
// Run:
//
//	go run ./
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/clivern/swarm"
	"github.com/google/uuid"
)

const (
	repoURL  = "https://github.com/octocat/Hello-World.git"
	deadline = 90 * time.Second
	sleepFor = 10 * time.Minute
)

const initBash = `
echo "init sleep start" >&2
sleep 600
echo "init sleep done" >&2
`

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	id := uuid.NewString()
	fmt.Printf("job id=%s\ndeadline=%s init sleep=%s\n", id, deadline, sleepFor)

	start := time.Now()
	_, err := swarm.Run(ctx, swarm.RunRequest{
		WorkDir:     "/tmp/basement",
		ID:          id,
		RepoURL:     repoURL,
		Prompt:      "Do not edit files. The init script should be killed before you start.",
		PIModel:     "openrouter/anthropic/claude-sonnet-4.5",
		ProxyURL:    "http://host.docker.internal:8080/api",
		DockerImage: "clivern/swarm:v0.8.5",
		Container: swarm.Container{
			InitBash: initBash,
			Memory:   "2g",
			CPUs:     "1",
		},
		Cleanup: true,
	})
	elapsed := time.Since(start)

	if err == nil {
		fmt.Fprintf(os.Stderr, "expected the deadline to stop the run after %s\n", elapsed.Round(time.Second))
		os.Exit(1)
	}
	if !strings.Contains(err.Error(), "init sleep start") {
		fmt.Fprintf(os.Stderr, "deadline fired before init started (%s): %v\n", elapsed.Round(time.Second), err)
		os.Exit(1)
	}
	if elapsed >= sleepFor {
		fmt.Fprintf(os.Stderr, "waited for the full init sleep (%s): %v\n", elapsed.Round(time.Second), err)
		os.Exit(1)
	}
	if !containerGone(id) {
		fmt.Fprintf(os.Stderr, "container %s still exists after %s: %v\n", id, elapsed.Round(time.Second), err)
		os.Exit(1)
	}

	fmt.Printf("init sleep killed after %s\n", elapsed.Round(time.Second))
	fmt.Printf("error: %v\n", err)
}

func containerGone(id string) bool {
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if exec.Command("docker", "inspect", id).Run() != nil {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}
