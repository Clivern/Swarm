// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunDocker(ctx context.Context, p DockerParams) error {
	repo, err := filepath.Abs(p.repoDir)
	if err != nil {
		return fmt.Errorf("resolve repo mount: %w", err)
	}

	out, err := filepath.Abs(p.outDir)
	if err != nil {
		return fmt.Errorf("resolve out mount: %w", err)
	}

	c := p.container
	env := []string{
		fmt.Sprintf("OPENROUTER_API_KEY=%s", p.openRouterAPIKey),
		fmt.Sprintf("PROMPT=%s", p.prompt),
		fmt.Sprintf("PI_MODEL=%s", p.piModel),
	}

	if p.initScript != "" {
		env = append(env, fmt.Sprintf("INIT_SCRIPT=%s", p.initScript))
	}

	args := []string{
		"run", "--rm",
		"--memory", c.Memory,
		"--cpus", c.CPUs,
	}

	for _, e := range env {
		args = append(args, "-e", e)
	}

	args = append(args,
		"-v", fmt.Sprintf("%s:/repo", repo),
		"-v", fmt.Sprintf("%s:/out", out),
		p.image,
	)

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

func WriteInitBash(outDir, bash string) error {
	if strings.TrimSpace(bash) == "" {
		return nil
	}

	path := filepath.Join(outDir, "init.sh")
	body := "#!/usr/bin/env bash\nset -euo pipefail\n" + bash + "\n"
	return os.WriteFile(path, []byte(body), 0o755)
}

func ResolveRepoInitScript(repoDir, rel string) (string, error) {
	rel = strings.TrimSpace(rel)
	if rel == "" {
		return "", nil
	}

	if filepath.IsAbs(rel) || strings.Contains(rel, "..") {
		return "", fmt.Errorf("Container.InitScript must be a path inside the repo")
	}

	rel = filepath.ToSlash(rel)
	host := filepath.Join(repoDir, filepath.FromSlash(rel))
	st, err := os.Stat(host)

	if err != nil {
		return "", fmt.Errorf("Container.InitScript: %w", err)
	}
	if st.IsDir() {
		return "", fmt.Errorf("Container.InitScript must be a file")
	}

	return rel, nil
}
