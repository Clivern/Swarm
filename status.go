// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	git "github.com/go-git/go-git/v5"
)

// ChangedFilesInRepo lists every path with a non-clean git status, sorted by path.
func ChangedFilesInRepo(repoDir string) ([]ChangedFile, error) {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("open repo: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("status: %w", err)
	}

	files := make([]ChangedFile, 0, len(status))
	for path, fs := range status {
		letter := FileStatusLetter(fs)
		if letter == "" {
			continue
		}

		var content string
		if letter != "D" {
			data, err := os.ReadFile(filepath.Join(repoDir, filepath.FromSlash(path)))
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			content = string(data)
		}
		files = append(files, ChangedFile{Path: path, Status: letter, Content: content})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// FileStatusLetter returns the porcelain status code, staging first. "" means unmodified.
func FileStatusLetter(fs *git.FileStatus) string {
	code := fs.Worktree
	if fs.Staging != git.Unmodified {
		code = fs.Staging
	}
	if code == git.Unmodified {
		return ""
	}
	return string(rune(code))
}
