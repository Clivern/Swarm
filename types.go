// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

// GitCloneAuth supplies credentials for cloning private repositories.
type GitCloneAuth struct {
	Token             string
	Username          string
	SSHPrivateKeyPath string
}

// RunRequest configures a Pi Docker run against a cloned repository.
type RunRequest struct {
	WorkDir      string
	ID           string
	RepoURL      string
	GitCloneAuth GitCloneAuth
	Prompt       string
	PIModel      string
	ProxyKey     string
	// ProxyURL is Pi's OpenRouter base URL.
	// Example: http://host.docker.internal:8080/api
	ProxyURL    string
	DockerImage string
	Container   Container
	Cleanup     bool
}

// Container configures init and docker run limits.
type Container struct {
	InitScript string
	InitBash   string
	Memory     string
	CPUs       string
}

// ChangedFile describes one path touched by the agent.
type ChangedFile struct {
	Path    string
	Status  string
	Content string
}

// Result holds artifacts from a successful run.
type Result struct {
	Patch        string
	Summary      string
	TotalTokens  int
	ChangedFiles []ChangedFile
	RepoDir      string
	OutDir       string
}

type DockerParams struct {
	Image      string
	ProxyKey   string
	ProxyURL   string
	Prompt     string
	PIModel    string
	RepoDir    string
	OutDir     string
	Container  Container
	InitScript string
}
