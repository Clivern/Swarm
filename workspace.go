// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"fmt"
	"os"
	"path/filepath"
)

// SetupWorkspace resolves WorkDir/ID/{repo,out} and creates the job directories.
func SetupWorkspace(workDir, id string) (repoDir, outDir string, err error) {
	root, err := filepath.Abs(filepath.Join(workDir, id))

	if err != nil {
		return "", "", fmt.Errorf("resolve WorkDir: %w", err)
	}

	repoDir = filepath.Join(root, "repo")
	outDir = filepath.Join(root, "out")

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create out dir: %w", err)
	}

	return repoDir, outDir, nil
}

// RemoveRepoDir deletes WorkDir/ID (repo and out). Other IDs under WorkDir are untouched.
func RemoveRepoDir(workDir, id string) error {
	if err := os.RemoveAll(filepath.Join(workDir, id)); err != nil {
		return fmt.Errorf("remove job dir: %w", err)
	}

	return nil
}
