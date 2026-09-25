## Swarm

Go library that clones a git repo, runs [Pi](https://pi.dev/) headless in Docker against it, and returns the unified diff.

### Prerequisites

- Go 1.25+
- Docker (running)
- [OpenRouter](https://openrouter.ai/) API key

Build the Pi image once from this repo:

```bash
docker build -t swarm-pi:local .
```

### Install

```bash
go get github.com/clivern/swarm
```

### Usage

```go
result, err := swarm.Run(ctx, swarm.RunRequest{
    WorkDir:          ".work",
    ID:               "550e8400-e29b-41d4-a716-446655440000",
    RepoURL:          "https://github.com/Clivern/Diffsay.git",
    Prompt:           "Add Gemfile.lock to LOW_PRIORITY_FILES in src/diffsay/const.py",
    PIModel:          "openrouter/anthropic/claude-sonnet-4.5",
    OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
    DockerImage:      "swarm-pi:local",
})

// result.Patch, result.Summary, result.RepoDir, result.OutDir
```

Each job uses `{WorkDir}/{ID}/repo` (clone) and `{WorkDir}/{ID}/out` (`patch.diff`, `summary.txt`). Reusing the same `ID` skips re-clone if `repo/.git` already exists.

### Private repositories

HTTPS:

```go
GitCloneAuth: swarm.GitCloneAuth{
    Token: os.Getenv("GITHUB_TOKEN"),
},
RepoURL: "https://github.com/org/private.git",
```

SSH:

```go
GitCloneAuth: swarm.GitCloneAuth{
    SSHPrivateKeyPath: os.ExpandEnv("$HOME/.ssh/id_ed25519"),
},
RepoURL: "git@github.com:org/private.git",
```
