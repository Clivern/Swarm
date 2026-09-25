// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
)

// WorkspacePaths returns repo and out directories for WorkDir/ID.
func WorkspacePaths(workDir, id string) (repoDir, outDir string, err error) {
	base, err := filepath.Abs(workDir)
	if err != nil {
		return "", "", fmt.Errorf("resolve WorkDir: %w", err)
	}

	root := filepath.Join(base, id)
	repoDir = filepath.Join(root, "repo")
	outDir = filepath.Join(root, "out")

	return repoDir, outDir, nil
}

// RemoveJobDir deletes WorkDir/ID (repo and out). Other IDs under WorkDir are untouched.
func RemoveJobDir(workDir, id string) error {
	repoDir, _, err := WorkspacePaths(workDir, id)
	if err != nil {
		return err
	}
	jobDir := filepath.Dir(repoDir)
	if err := os.RemoveAll(jobDir); err != nil {
		return fmt.Errorf("remove job dir: %w", err)
	}
	return nil
}

// EnsureClone clones repoURL into dest when dest is not already a git repository.
func EnsureClone(ctx context.Context, repoURL, dest string, auth GitCloneAuth) error {
	gitDir := filepath.Join(dest, ".git")
	if st, err := os.Stat(gitDir); err == nil && st.IsDir() {
		return nil
	}

	repoURL = strings.TrimSpace(repoURL)
	if repoURL == "" {
		return fmt.Errorf("repo URL is empty")
	}

	cloneAuth, err := cloneAuth(repoURL, auth)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("create workspace dir: %w", err)
	}

	_, err = git.PlainCloneContext(ctx, dest, false, &git.CloneOptions{
		URL:   repoURL,
		Depth: 1,
		Auth:  cloneAuth,
	})
	if err != nil {
		return fmt.Errorf("clone: %w", err)
	}

	return nil
}
