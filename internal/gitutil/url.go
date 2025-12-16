package gitutil

import (
	"fmt"
	"strings"
)

func CommitURLBase(ciServerURL, ciProjectPath string) (string, error) {
	// CI-friendly: use GitLab-provided values if present
	if ciServerURL != "" && ciProjectPath != "" {
		return strings.TrimRight(ciServerURL, "/") + "/" + strings.TrimLeft(ciProjectPath, "/") + "/-/commit", nil
	}

	// Local fallback: parse origin URL
	originURL, err := runGit("remote", "get-url", "origin")
	if err != nil {
		return "", fmt.Errorf("get origin url: %w", err)
	}
	originURL = strings.TrimSpace(originURL)
	originURL = strings.TrimSuffix(originURL, ".git")

	// https://gitlab.example.com/group/project
	if strings.HasPrefix(originURL, "http://") || strings.HasPrefix(originURL, "https://") {
		// strip scheme
		hostAndPath := originURL
		if i := strings.Index(hostAndPath, "://"); i >= 0 {
			hostAndPath = hostAndPath[i+3:]
		}
		// host/path
		slash := strings.Index(hostAndPath, "/")
		if slash < 0 {
			return "", fmt.Errorf("unexpected origin url: %q", originURL)
		}
		host := hostAndPath[:slash]
		path := hostAndPath[slash+1:]
		return "https://" + host + "/" + path + "/-/commit", nil
	}

	// git@gitlab.example.com:group/project
	if strings.HasPrefix(originURL, "git@") {
		rest := strings.TrimPrefix(originURL, "git@")
		colon := strings.Index(rest, ":")
		if colon < 0 {
			return "", fmt.Errorf("unexpected origin url: %q", originURL)
		}
		host := rest[:colon]
		path := rest[colon+1:]
		return "https://" + host + "/" + path + "/-/commit", nil
	}

	return "", fmt.Errorf("unsupported origin url: %q", originURL)
}