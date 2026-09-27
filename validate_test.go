// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validRunRequest() RunRequest {
	return RunRequest{
		WorkDir:     "/tmp/basement",
		ID:          "job-1",
		RepoURL:     "https://github.com/example/repo.git",
		Prompt:      "do something",
		PIModel:     "openrouter/anthropic/claude-sonnet-4.5",
		ProxyURL:    "https://openrouter.ai/api/v1",
		DockerImage: "swarm-pi:local",
		Container: Container{
			Memory: "2g",
			CPUs:   "1",
		},
	}
}

func TestUnitRunRequestValidate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		assert.NoError(t, validRunRequest().validate())
	})

	t.Run("missing repository URL", func(t *testing.T) {
		req := validRunRequest()
		req.RepoURL = ""
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Repository URL is required")
	})

	t.Run("missing prompt", func(t *testing.T) {
		req := validRunRequest()
		req.Prompt = "  "
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Prompt is required")
	})

	t.Run("missing PI model", func(t *testing.T) {
		req := validRunRequest()
		req.PIModel = ""
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "PI Model is required")
	})

	t.Run("missing ID", func(t *testing.T) {
		req := validRunRequest()
		req.ID = ""
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "ID is required")
	})

	t.Run("missing proxy URL", func(t *testing.T) {
		req := validRunRequest()
		req.ProxyURL = "  "
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "ProxyURL is required")
	})

	t.Run("missing Docker image", func(t *testing.T) {
		req := validRunRequest()
		req.DockerImage = ""
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Docker Image is required")
	})
}
