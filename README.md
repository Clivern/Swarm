## Swarm

Swarm runs coding agents on real codebases: clone any git repository (public or private), execute [Pi](https://pi.dev/) headless in Docker, and return a unified diff - so you can review, apply, or open a PR from automated tasks.


### Prerequisites

- Go 1.25+
- Docker (running)
- [OpenRouter](https://openrouter.ai/) API key

Build the Pi image once from this repo:

```bash
docker build -t swarm:v0.1.0 .
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
    Prompt:           "Add Gemfile.lock to LOW_PRIORITY_FILES",
    PIModel:          "openrouter/anthropic/claude-sonnet-4.5",
    OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
    DockerImage:      "swarm:v0.1.0",
})

// result.Patch, result.Summary, result.RepoDir, result.OutDir
```


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
