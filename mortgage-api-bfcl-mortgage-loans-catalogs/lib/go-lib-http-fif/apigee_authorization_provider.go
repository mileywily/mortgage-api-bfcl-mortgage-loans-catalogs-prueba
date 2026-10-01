package http_fif

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const apikeyTokenExpiryBuffer = 5 * time.Second

type apikeyAuthorizationProvider struct {
	clientID     string
	clientSecret string
	authURL      string
	httpClient   *http.Client

	token     string
	expiresAt time.Time
	mu        sync.Mutex
}

func NewApiKeyAuthorizationProvider(clientID, clientSecret, authURL string, httpClient *http.Client) AuthorizationProvider {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &apikeyAuthorizationProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		authURL:      authURL,
		httpClient:   httpClient,
	}
}

func (p *apikeyAuthorizationProvider) GetAuthorization() (out string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	defer func() {
		out = "Bearer " + out
	}()

	if time.Now().Before(p.expiresAt) {
		return p.token, nil
	}

	return p.refreshToken()
}

func (p *apikeyAuthorizationProvider) refreshToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", p.authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	clientIDAndSecret := fmt.Sprintf("%s:%s", p.clientID, p.clientSecret)
	authToken := base64.StdEncoding.EncodeToString([]byte(clientIDAndSecret))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("x-api-key", p.clientID)
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authToken))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch token")
	}

	var response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	expiresIn, err := strconv.ParseInt(response.ExpiresIn, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid expires_in value in response: %q", response.ExpiresIn)
	}

	p.token = response.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(expiresIn)*time.Second - apikeyTokenExpiryBuffer)
	return p.token, nil
}
