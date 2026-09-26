// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Run clones the repo if needed, runs the Pi container, and returns the patch.
func Run(ctx context.Context, req RunRequest) (*Result, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	repoDir, outDir, err := SetupWorkspace(req.WorkDir, req.ID)
	if err != nil {
		return nil, err
	}

	if err := EnsureClone(ctx, req.RepoURL, repoDir, req.GitCloneAuth); err != nil {
		return nil, err
	}

	initScript, err := ResolveRepoInitScript(repoDir, req.Container.InitScript)
	if err != nil {
		return nil, err
	}

	if err := WriteInitBash(outDir, req.Container.InitBash); err != nil {
		return nil, fmt.Errorf("write inline init: %w", err)
	}

	if err := RunDocker(ctx, DockerParams{
		Image:      req.DockerImage,
		ProxyKey:   req.ProxyKey,
		ProxyURL:   strings.TrimSpace(req.ProxyURL),
		Prompt:     req.Prompt,
		PIModel:    req.PIModel,
		RepoDir:    repoDir,
		OutDir:     outDir,
		Container:  req.Container,
		InitScript: initScript,
	}); err != nil {
		return nil, err
	}

	patch, err := os.ReadFile(filepath.Join(outDir, "patch.diff"))
	if err != nil {
		return nil, fmt.Errorf("read patch.diff: %w", err)
	}

	summary, totalTokens, err := ReadOut(outDir)
	if err != nil {
		return nil, err
	}

	changed, err := ChangedFilesInRepo(repoDir)
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
		if err := RemoveRepoDir(req.WorkDir, req.ID); err != nil {
			return nil, err
		}
	}

	return result, nil
}
