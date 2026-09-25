package swarm

import (
	"fmt"
	"strings"
)

func (req RunRequest) validate() error {
	if strings.TrimSpace(req.RepoURL) == "" {
		return fmt.Errorf("Repository URL is required")
	}
	if strings.TrimSpace(req.Prompt) == "" {
		return fmt.Errorf("Prompt is required")
	}
	if strings.TrimSpace(req.PIModel) == "" {
		return fmt.Errorf("PI Model is required")
	}
	if strings.TrimSpace(req.OpenRouterAPIKey) == "" {
		return fmt.Errorf("OpenRouterAPIKey is required")
	}
	if strings.TrimSpace(req.DockerImage) == "" {
		return fmt.Errorf("Docker Image is required")
	}
	return nil
}
