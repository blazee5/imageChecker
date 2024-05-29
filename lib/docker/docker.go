package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	defaultRegistry = "index.docker.io"
	defaultTag      = "latest"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func ParseDockerImage(image string) (string, string, string) {
	var registry, repository, tag string

	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		image = parts[0]
		tag = parts[1]
	} else {
		tag = defaultTag
	}

	slashParts := strings.Split(image, "/")

	switch {
	case len(slashParts) == 1:
		registry = defaultRegistry
		repository = "library/" + image
	case len(slashParts) == 2:
		registry = defaultRegistry
		repository = image
	default:
		registry = slashParts[0]
		repository = strings.Join(slashParts[1:], "/")
	}

	return registry, repository, tag
}

func AuthDockerHub(ctx context.Context, repository, username, password string) (string, error) {
	scope := fmt.Sprintf("repository:%s:pull", repository)
	url := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=%s", scope)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	if username != "" && password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("error while do auth request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error while get auth token, status code: %d", resp.StatusCode)
	}

	var response AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func AuthOtherRegistry(username, password string) (string, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))

	return auth, nil
}

func CheckImage(ctx context.Context, registry, repository, tag, username, password string) (bool, error) {
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", registry, repository, tag)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false, err
	}

	if registry == defaultRegistry {
		token, err := AuthDockerHub(ctx, repository, username, password)
		if err != nil {
			return false, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		token, err := AuthOtherRegistry(username, password)
		if err != nil {
			return false, err
		}

		req.Header.Set("Authorization", "Basic "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}
