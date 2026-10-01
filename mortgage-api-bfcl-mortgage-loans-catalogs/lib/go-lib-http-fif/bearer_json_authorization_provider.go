package http_fif

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

type jsonAuthorizationProvider struct {
	tokenEndpoint string
	client        *http.Client
	credentials   map[string]string

	token     string
	expiresAt time.Time
	mu        sync.Mutex
}

func NewJSONAuthorizationProvider(tokenEndpoint string, client *http.Client, credentials map[string]string) AuthorizationProvider {
	if client == nil {
		client = http.DefaultClient
	}

	return &jsonAuthorizationProvider{
		tokenEndpoint: tokenEndpoint,
		client:        client,
		credentials:   credentials,
	}
}

func (p *jsonAuthorizationProvider) GetAuthorization() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if time.Now().Before(p.expiresAt) {
		return "Bearer " + p.token, nil
	}

	return p.refreshToken()
}

func (p *jsonAuthorizationProvider) refreshToken() (string, error) {
	payload, err := json.Marshal(p.credentials)
	if err != nil {
		return "", errors.New("error encoding JSON payload")
	}

	req, err := http.NewRequest("POST", p.tokenEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cache-control", "no-cache")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch token")
	}

	var response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	p.token = response.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)

	return "Bearer " + p.token, nil
}
