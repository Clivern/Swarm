// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

func cloneAuth(repoURL string, auth GitCloneAuth) (transport.AuthMethod, error) {
	repoURL = strings.TrimSpace(repoURL)
	if repoURL == "" {
		return nil, fmt.Errorf("repo URL is empty")
	}

	if strings.HasPrefix(repoURL, "git@") || strings.HasPrefix(repoURL, "ssh:") {
		key := strings.TrimSpace(auth.SSHPrivateKeyPath)
		if key == "" && strings.TrimSpace(auth.Token) != "" {
			return nil, fmt.Errorf("HTTPS token was set but RepoURL is SSH; use an https:// URL or set SSHPrivateKeyPath")
		}
		if key == "" {
			return nil, nil
		}
		absKey, err := filepath.Abs(key)
		if err != nil {
			return nil, fmt.Errorf("resolve SSHPrivateKeyPath: %w", err)
		}
		pub, err := gitssh.NewPublicKeysFromFile("git", absKey, "")
		if err != nil {
			return nil, fmt.Errorf("load SSH key: %w", err)
		}
		pub.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		return pub, nil
	}

	token := strings.TrimSpace(auth.Token)
	if token == "" {
		return nil, nil
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, fmt.Errorf("parse RepoURL: %w", err)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return nil, fmt.Errorf("Git token auth requires an https:// RepoURL (got scheme %q)", u.Scheme)
	}

	user := strings.TrimSpace(auth.Username)
	if user == "" {
		user = "x-access-token"
	}
	return &githttp.BasicAuth{Username: user, Password: token}, nil
}
