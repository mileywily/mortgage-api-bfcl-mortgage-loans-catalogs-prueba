package http_fif

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const bearerTokenExpiryBuffer = 5 * time.Second

type bearerAuthorizationProvider struct {
	clientID     string
	clientSecret string
	authURL      string
	httpClient   *http.Client

	token     string
	expiresAt time.Time
	mu        sync.Mutex
}

func NewBearerAuthorizationProvider(clientID, clientSecret, authURL string, httpClient *http.Client) AuthorizationProvider {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &bearerAuthorizationProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		authURL:      authURL,
		httpClient:   httpClient,
	}
}

func (p *bearerAuthorizationProvider) GetAuthorization() (out string, err error) {
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

func (p *bearerAuthorizationProvider) refreshToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)

	req, err := http.NewRequest("POST", p.authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		ExpiresIn   int64  `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	p.token = response.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(response.ExpiresIn)*time.Second - bearerTokenExpiryBuffer)
	return p.token, nil
}
