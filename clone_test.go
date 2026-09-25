// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"testing"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
)

func TestUnitCloneAuth(t *testing.T) {
	t.Run("public HTTPS", func(t *testing.T) {
		auth, err := cloneAuth("https://github.com/org/repo.git", GitCloneAuth{})
		assert.NoError(t, err)
		assert.Nil(t, auth)
	})

	t.Run("HTTPS token default username", func(t *testing.T) {
		auth, err := cloneAuth("https://github.com/org/private.git", GitCloneAuth{Token: "secret"})
		assert.NoError(t, err)

		basic, ok := auth.(*githttp.BasicAuth)
		assert.True(t, ok)
		assert.Equal(t, "x-access-token", basic.Username)
		assert.Equal(t, "secret", basic.Password)
	})

	t.Run("HTTPS token custom username", func(t *testing.T) {
		auth, err := cloneAuth("https://gitlab.com/g/r.git", GitCloneAuth{
			Token:    "tok",
			Username: "oauth2",
		})
		assert.NoError(t, err)

		basic, ok := auth.(*githttp.BasicAuth)
		assert.True(t, ok)
		assert.Equal(t, "oauth2", basic.Username)
		assert.Equal(t, "tok", basic.Password)
	})

	t.Run("SSH URL with HTTPS token", func(t *testing.T) {
		_, err := cloneAuth("git@github.com:org/repo.git", GitCloneAuth{Token: "x"})
		assert.Error(t, err)
	})

	t.Run("SSH URL without key", func(t *testing.T) {
		auth, err := cloneAuth("git@github.com:org/repo.git", GitCloneAuth{})
		assert.NoError(t, err)
		assert.Nil(t, auth)
	})

	t.Run("empty repository URL", func(t *testing.T) {
		_, err := cloneAuth("  ", GitCloneAuth{})
		assert.Error(t, err)
	})
}
