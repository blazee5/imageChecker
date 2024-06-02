package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
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

func AuthDockerHub(ctx context.Context, repository, username, password string, timeout time.Duration) (string, error) {
	scope := fmt.Sprintf("repository:%s:pull", repository)
	url := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=%s", scope)

	client := resty.New().SetTimeout(timeout)
	req := client.R().SetContext(ctx)

	if username != "" && password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
		req.SetHeader("Authorization", "Basic "+auth)
	}

	resp, err := req.Get(url)
	if err != nil {
		return "", fmt.Errorf("error while doing auth request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("error while getting auth token, status code: %d", resp.StatusCode())
	}

	var response AuthResponse
	err = json.Unmarshal(resp.Body(), &response)

	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func AuthOtherRegistry(username, password string) (string, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	return auth, nil
}

func CheckImage(ctx context.Context, registry, repository, tag, username, password string, timeout time.Duration) (bool, error) {
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", registry, repository, tag)

	client := resty.New().SetTimeout(timeout)
	req := client.R().SetContext(ctx).SetHeader("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	if registry == defaultRegistry {
		token, err := AuthDockerHub(ctx, repository, username, password, timeout)
		if err != nil {
			return false, err
		}

		req.SetHeader("Authorization", "Bearer "+token)
	} else {
		token, err := AuthOtherRegistry(username, password)
		if err != nil {
			return false, err
		}

		req.SetHeader("Authorization", "Basic "+token)
	}

	resp, err := req.Head(url)
	if err != nil {
		return false, err
	}

	if resp.StatusCode() != http.StatusOK {
		return false, nil
	}

	return true, nil
}
