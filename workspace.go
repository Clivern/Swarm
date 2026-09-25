package swarm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WorkspacePaths returns repo and out directories for WorkDir/ID.
func WorkspacePaths(workDir, id string) (repoDir, outDir string, err error) {
	base, err := filepath.Abs(workDir)

	if err != nil {
		return "", "", fmt.Errorf("resolve WorkDir: %w", err)
	}

	root := filepath.Join(base, id)
	repoDir = filepath.Join(root, "repo")
	outDir = filepath.Join(root, "out")

	return repoDir, outDir, nil
}

// EnsureClone clones repoURL into dest when dest is not already a git repository.
func EnsureClone(ctx context.Context, repoURL, dest string, auth GitCloneAuth) error {
	gitDir := filepath.Join(dest, ".git")
	if st, err := os.Stat(gitDir); err == nil && st.IsDir() {
		return nil
	}

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("create repo dir: %w", err)
	}

	cloneURL, sshKey, err := prepareCloneURL(repoURL, auth)
	if err != nil {
		return err
	}

	args := []string{"clone", "--depth", "1"}
	if sshKey != "" {
		absKey, err := filepath.Abs(sshKey)
		if err != nil {
			return fmt.Errorf("resolve SSHPrivateKeyPath: %w", err)
		}
		sshCmd := fmt.Sprintf("ssh -i %q -o BatchMode=yes -o StrictHostKeyChecking=accept-new", absKey)
		args = append([]string{"-c", "core.sshCommand=" + sshCmd}, args...)
	}

	args = append(args, cloneURL, dest)
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.CombinedOutput()

	if err != nil {
		msg := strings.TrimSpace(string(out))
		msg = redactCloneSecrets(msg, auth)
		if len(msg) > 4096 {
			msg = msg[:4096] + "..."
		}
		return fmt.Errorf("git clone: %w: %s", err, msg)
	}

	return nil
}
