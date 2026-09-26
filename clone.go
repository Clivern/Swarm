// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

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
