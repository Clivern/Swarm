// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

// EnsureClone shallow-clones repoURL into dest unless dest is already a git repository.
func EnsureClone(ctx context.Context, repoURL, dest string, auth GitCloneAuth) error {
	if st, err := os.Stat(filepath.Join(dest, ".git")); err == nil && st.IsDir() {
		return nil
	}

	method, err := CloneAuth(repoURL, auth)
	if err != nil {
		return err
	}

	if _, err := git.PlainCloneContext(ctx, dest, false, &git.CloneOptions{
		URL:   repoURL,
		Depth: 1,
		Auth:  method,
	}); err != nil {
		return fmt.Errorf("clone: %w", err)
	}

	return nil
}

// CloneAuth returns SSH key auth for git@/ssh: URLs, token auth for HTTPS, or nil for public repos.
func CloneAuth(repoURL string, auth GitCloneAuth) (transport.AuthMethod, error) {
	if strings.HasPrefix(repoURL, "git@") || strings.HasPrefix(repoURL, "ssh:") {
		if auth.SSHPrivateKeyPath == "" {
			return nil, nil
		}
		pub, err := gitssh.NewPublicKeysFromFile("git", auth.SSHPrivateKeyPath, "")
		if err != nil {
			return nil, fmt.Errorf("load SSH key: %w", err)
		}
		pub.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		return pub, nil
	}

	if auth.Token == "" {
		return nil, nil
	}

	user := auth.Username
	if user == "" {
		user = "x-access-token"
	}

	return &githttp.BasicAuth{Username: user, Password: auth.Token}, nil
}
