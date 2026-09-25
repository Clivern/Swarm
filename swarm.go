// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

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

	summary, totalTokens, err := readOut(outDir)
	if err != nil {
		return nil, err
	}

	changed, err := changedFilesInRepo(repoDir)
	if err != nil {
		return nil, fmt.Errorf("list changed files: %w", err)
	}

	result := &Result{
		Patch:        string(patch),
		Summary:      summary,
		TotalTokens:  totalTokens,
		ChangedFiles: changed,
		RepoDir:      repoDir,
		OutDir:       outDir,
	}

	if req.Cleanup {
		if err := RemoveJobDir(req.WorkDir, req.ID); err != nil {
			return nil, err
		}
	}

	return result, nil
}
