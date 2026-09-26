// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"fmt"
	"strings"
)

func (req RunRequest) validate() error {
	required := []struct{ name, value string }{
		{"Repository URL", req.RepoURL},
		{"Prompt", req.Prompt},
		{"PI Model", req.PIModel},
		{"OpenRouterAPIKey", req.OpenRouterAPIKey},
		{"Docker Image", req.DockerImage},
		{"Container.Memory", req.Container.Memory},
		{"Container.CPUs", req.Container.CPUs},
	}
	for _, f := range required {
		if strings.TrimSpace(f.value) == "" {
			return fmt.Errorf("%s is required", f.name)
		}
	}
	return nil
}
