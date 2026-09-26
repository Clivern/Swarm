// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitSetupWorkspace(t *testing.T) {
	root := t.TempDir()
	repo, out, err := SetupWorkspace(root, "abc-123")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "abc-123", "repo"), repo)
	assert.Equal(t, filepath.Join(root, "abc-123", "out"), out)
	assert.DirExists(t, out)

	assert.NoError(t, RemoveRepoDir(root, "abc-123"))
	assert.NoDirExists(t, filepath.Join(root, "abc-123"))
}
