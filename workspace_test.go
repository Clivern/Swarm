// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitWorkspace(t *testing.T) {
	t.Run("WorkspacePaths", func(t *testing.T) {
		root := t.TempDir()
		repo, out, err := WorkspacePaths(root, "abc-123")
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(root, "abc-123", "repo"), repo)
		assert.Equal(t, filepath.Join(root, "abc-123", "out"), out)
	})

	t.Run("EnsureClone skips when git present", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "repo")
		assert.NoError(t, os.MkdirAll(filepath.Join(dest, ".git"), 0o755))

		err := EnsureClone(context.Background(), "https://example.com/nope.git", dest, GitCloneAuth{})
		assert.NoError(t, err)
	})

	t.Run("EnsureClone public repository", func(t *testing.T) {
		if testing.Short() {
			t.Skip("network clone")
		}
		dest := filepath.Join(t.TempDir(), "repo")
		err := EnsureClone(
			context.Background(),
			"https://github.com/octocat/Hello-World.git",
			dest,
			GitCloneAuth{},
		)
		assert.NoError(t, err)
		_, err = os.Stat(filepath.Join(dest, ".git"))
		assert.NoError(t, err)
	})
}
