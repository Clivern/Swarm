// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	git "github.com/go-git/go-git/v5"
)

const maxChangedFileBytes = 1 << 20 // 1 MiB

func changedFilesInRepo(repoDir string) ([]ChangedFile, error) {
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
		if fs == nil || (fs.Staging == git.Unmodified && fs.Worktree == git.Unmodified) {
			continue
		}
		letter := fileStatusLetter(fs)
		var content string
		if letter != "D" {
			content, err = readWorktreeFile(repoDir, path)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", path, err)
			}
		}
		files = append(files, ChangedFile{
			Path:    path,
			Status:  letter,
			Content: content,
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func fileStatusLetter(fs *git.FileStatus) string {
	code := fs.Worktree
	if fs.Staging != git.Unmodified {
		code = fs.Staging
	}
	if code == git.Unmodified {
		return ""
	}
	return string([]byte{byte(code)})
}

func readWorktreeFile(repoDir, path string) (string, error) {
	abs := filepath.Join(repoDir, filepath.FromSlash(path))
	info, err := os.Lstat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", nil
	}
	f, err := os.Open(abs)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxChangedFileBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxChangedFileBytes {
		return string(data[:maxChangedFileBytes]) + "\n… truncated", nil
	}
	return string(data), nil
}
