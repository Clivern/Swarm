// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
)

func TestUnitEnsureClone(t *testing.T) {
	t.Run("skips when git present", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "repo")
		assert.NoError(t, os.MkdirAll(filepath.Join(dest, ".git"), 0o755))

		err := EnsureClone(context.Background(), "https://example.com/nope.git", dest, GitCloneAuth{})
		assert.NoError(t, err)
	})

	t.Run("public repository", func(t *testing.T) {
		if testing.Short() {
			t.Skip("network clone")
		}
		dest := filepath.Join(t.TempDir(), "repo")
		err := EnsureClone(context.Background(), "https://github.com/octocat/Hello-World.git", dest, GitCloneAuth{})
		assert.NoError(t, err)
		assert.DirExists(t, filepath.Join(dest, ".git"))
	})
}

func TestUnitCloneAuth(t *testing.T) {
	t.Run("public HTTPS", func(t *testing.T) {
		auth, err := CloneAuth("https://github.com/org/repo.git", GitCloneAuth{})
		assert.NoError(t, err)
		assert.Nil(t, auth)
	})

	t.Run("HTTPS token default username", func(t *testing.T) {
		auth, err := CloneAuth("https://github.com/org/private.git", GitCloneAuth{Token: "secret"})
		assert.NoError(t, err)

		basic, ok := auth.(*githttp.BasicAuth)
		assert.True(t, ok)
		assert.Equal(t, "x-access-token", basic.Username)
		assert.Equal(t, "secret", basic.Password)
	})

	t.Run("HTTPS token custom username", func(t *testing.T) {
		auth, err := CloneAuth("https://gitlab.com/g/r.git", GitCloneAuth{
			Token:    "tok",
			Username: "oauth2",
		})
		assert.NoError(t, err)

		basic, ok := auth.(*githttp.BasicAuth)
		assert.True(t, ok)
		assert.Equal(t, "oauth2", basic.Username)
		assert.Equal(t, "tok", basic.Password)
	})

	t.Run("SSH URL without key", func(t *testing.T) {
		auth, err := CloneAuth("git@github.com:org/repo.git", GitCloneAuth{})
		assert.NoError(t, err)
		assert.Nil(t, auth)
	})
}
