package internal

import (
	"errors"
	"os/exec"
	"strings"
)

func GetCurrentGitRef() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetCurrentGitProjectName() (string, string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", "", err
	}

	remoteURL := strings.TrimSpace(string(output))

	switch {
	case strings.HasPrefix(remoteURL, "git@"):
		remoteURL = strings.TrimPrefix(remoteURL, "git@")
		remoteURL = strings.TrimSuffix(remoteURL, ".git")

		parts := strings.SplitN(remoteURL, ":", 2)
		if len(parts) != 2 {
			return "", "", errors.New("can not parse git remote SSH URL")
		}

		return "https://" + parts[0], parts[1], nil
	case strings.HasPrefix(remoteURL, "http://"), strings.HasPrefix(remoteURL, "https://"):
		remoteURL = strings.TrimPrefix(remoteURL, "http://")
		remoteURL = strings.TrimPrefix(remoteURL, "https://")

		remoteURL = strings.TrimSuffix(remoteURL, ".git")

		parts := strings.SplitN(remoteURL, "/", 2)
		if len(parts) != 2 {
			return "", "", errors.New("can not parse git remote HTTPS URL")
		}
		return "https://" + parts[0], parts[1], nil
	}

	return "", "", errors.New("can not parse git remote")
}
