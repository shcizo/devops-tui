package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsGitRepo checks if the current directory is a git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("getting current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// BranchExists checks if a branch exists
func BranchExists(name string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+name)
	err := cmd.Run()
	return err == nil
}

// CreateBranch creates a new branch and optionally checks it out
func CreateBranch(name string, checkout bool) error {
	if BranchExists(name) {
		return fmt.Errorf("branch '%s' already exists", name)
	}

	if checkout {
		// Create and checkout
		cmd := exec.Command("git", "checkout", "-b", name)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("creating branch: %s", strings.TrimSpace(string(output)))
		}
	} else {
		// Just create
		cmd := exec.Command("git", "branch", name)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("creating branch: %s", strings.TrimSpace(string(output)))
		}
	}

	return nil
}

// CheckoutBranch checks out an existing branch
func CheckoutBranch(name string) error {
	cmd := exec.Command("git", "checkout", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("checking out branch: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

// HasUncommittedChanges checks if there are uncommitted changes
func HasUncommittedChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}
