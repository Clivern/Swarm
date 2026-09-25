package swarm

import (
	"fmt"
	"net/url"
	"strings"
)

func prepareCloneURL(repoURL string, auth GitCloneAuth) (cloneURL string, sshKey string, err error) {
	repoURL = strings.TrimSpace(repoURL)

	if repoURL == "" {
		return "", "", fmt.Errorf("repo URL is empty")
	}

	if strings.HasPrefix(repoURL, "git@") || strings.HasPrefix(repoURL, "ssh:") {
		key := strings.TrimSpace(auth.SSHPrivateKeyPath)
		if key == "" && strings.TrimSpace(auth.Token) != "" {
			return "", "", fmt.Errorf("HTTPS token was set but RepoURL is SSH; use an https:// URL or set SSHPrivateKeyPath")
		}
		return repoURL, key, nil
	}

	token := strings.TrimSpace(auth.Token)
	if token == "" {
		return repoURL, "", nil
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", fmt.Errorf("parse RepoURL: %w", err)
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return "", "", fmt.Errorf("Git token auth requires an https:// RepoURL (got scheme %q)", u.Scheme)
	}

	user := strings.TrimSpace(auth.Username)

	u.User = url.UserPassword(user, token)
	return u.String(), "", nil
}

func redactCloneSecrets(msg string, auth GitCloneAuth) string {
	token := strings.TrimSpace(auth.Token)

	if token != "" {
		msg = strings.ReplaceAll(msg, token, "***")
	}

	return msg
}
