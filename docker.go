// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func runDocker(ctx context.Context, p dockerParams) error {
	repo, err := filepath.Abs(p.repoDir)
	if err != nil {
		return fmt.Errorf("resolve repo mount: %w", err)
	}

	out, err := filepath.Abs(p.outDir)
	if err != nil {
		return fmt.Errorf("resolve out mount: %w", err)
	}

	args := []string{
		"run", "--rm",
		"-e", "OPENROUTER_API_KEY=" + p.openRouterAPIKey,
		"-e", "PROMPT=" + p.prompt,
		"-e", "PI_MODEL=" + p.piModel,
		"-v", repo + ":/repo",
		"-v", out + ":/out",
		p.image,
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("docker run: %s", msg)
	}

	return nil
}
