package swarm

// GitCloneAuth supplies credentials for cloning private repositories.
// Use Token for HTTPS URLs; use SSHPrivateKeyPath for git@ or ssh:// URLs.
type GitCloneAuth struct {
	// Token is the HTTPS password (e.g. GitHub/GitLab personal access token).
	Token string
	// Username for HTTPS basic auth. Empty with Token defaults to x-access-token (GitHub).
	// GitLab deploy tokens often use gitlab-ci-token; OAuth-style clones may use oauth2.
	Username string
	// SSHPrivateKeyPath is passed to git as core.sshCommand for SSH remotes.
	SSHPrivateKeyPath string
}

// RunRequest configures a Pi Docker run against a cloned repository.
type RunRequest struct {
	// WorkDir is the parent directory; each job uses WorkDir/ID/.
	WorkDir string
	// ID uniquely names the workspace (e.g. a UUID).
	ID string
	// RepoURL is cloned into WorkDir/ID/repo when missing.
	RepoURL string
	// GitCloneAuth authenticates private repo clones (HTTPS token or SSH key).
	GitCloneAuth GitCloneAuth

	Prompt           string
	PIModel          string
	OpenRouterAPIKey string

	// DockerImage is the Pi image tag (e.g. swarm-pi:local).
	DockerImage string
}

// Result holds artifacts from a successful run.
type Result struct {
	Patch   string
	Summary string
	RepoDir string
	OutDir  string
}

type dockerParams struct {
	image            string
	openRouterAPIKey string
	prompt           string
	piModel          string
	repoDir          string
	outDir           string
}
