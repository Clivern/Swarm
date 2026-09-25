package swarm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Run clones the repo if needed, runs the Pi container, and returns the patch.
func Run(ctx context.Context, req RunRequest) (*Result, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	repoDir, outDir, err := WorkspacePaths(req.WorkDir, req.ID)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(repoDir), 0o755); err != nil {
		return nil, fmt.Errorf("create workspace parent: %w", err)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("create out dir: %w", err)
	}
	if err := EnsureClone(ctx, req.RepoURL, repoDir, req.GitCloneAuth); err != nil {
		return nil, err
	}

	if err := runDocker(ctx, dockerParams{
		image:            req.DockerImage,
		openRouterAPIKey: req.OpenRouterAPIKey,
		prompt:           req.Prompt,
		piModel:          req.PIModel,
		repoDir:          repoDir,
		outDir:           outDir,
	}); err != nil {
		return nil, err
	}

	patch, err := os.ReadFile(filepath.Join(outDir, "patch.diff"))
	if err != nil {
		return nil, fmt.Errorf("read patch.diff: %w", err)
	}

	summary, err := os.ReadFile(filepath.Join(outDir, "summary.txt"))
	if err != nil {
		return nil, fmt.Errorf("read summary.txt: %w", err)
	}

	return &Result{
		Patch:   string(patch),
		Summary: string(summary),
		RepoDir: repoDir,
		OutDir:  outDir,
	}, nil
}
