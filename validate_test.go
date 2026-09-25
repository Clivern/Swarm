// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validRunRequest() RunRequest {
	return RunRequest{
		WorkDir:          ".work",
		ID:               "job-1",
		RepoURL:          "https://github.com/example/repo.git",
		Prompt:           "do something",
		PIModel:          "openrouter/anthropic/claude-sonnet-4.5",
		OpenRouterAPIKey: "sk-or-test",
		DockerImage:      "swarm-pi:local",
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

	t.Run("missing OpenRouter API key", func(t *testing.T) {
		req := validRunRequest()
		req.OpenRouterAPIKey = ""
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "OpenRouterAPIKey is required")
	})

	t.Run("missing Docker image", func(t *testing.T) {
		req := validRunRequest()
		req.DockerImage = ""
		err := req.validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Docker Image is required")
	})
}
