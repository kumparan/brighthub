package brighthub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	Client interface {
		AddVideoToFolder(videoID, folderID string) error
		CreateVideo(req *CreateVideoRequest) (*CreateVideoResponse, error)
		IngestVideo(videoID string, req *IngestVideoRequest) (*IngestVideoResponse, error)
	}

	client struct {
		accessToken           string
		accessTokenAcquiredAt time.Time
		accountID             string
		clientID              string
		clientSecret          string
		httpClient            *http.Client
	}

	getAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
)

const authBaseURL = "https://oauth.brightcove.com/v4"

var defaultHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

// New :nodoc:
func New(clientID, clientSecret, accountID string, httpClient *http.Client) (Client, error) {
	c := &client{
		accountID:    accountID,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
	}
	if httpClient == nil {
		c.httpClient = defaultHTTPClient
	}

	_, err := c.getAccessToken()
	if err != nil {
		log.WithFields(log.Fields{
			"client_id":     clientID,
			"client_secret": clientSecret}).
			Error(err)
		return nil, err
	}
	return c, nil
}

func (c *client) getAccessToken() (string, error) {
	// Access Token only valid for 5 minutes. If > 5 minutes then get another token and update.
	// Since we cannot sure, therefore make a 1 minute buffer.
	if time.Since(c.accessTokenAcquiredAt).Minutes() <= 4 {
		return c.accessToken, nil
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/access_token?grant_type=client_credentials", authBaseURL), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret}).
			Error(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.clientID+":"+c.clientSecret)))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret}).
			Error(err)
		return "", err
	}
	defer resp.Body.Close()

	a := new(getAccessTokenResponse)
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		log.WithFields(log.Fields{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret}).
			Error(err)
		return "", err
	}
	c.accessToken = a.AccessToken
	c.accessTokenAcquiredAt = time.Now()

	return a.AccessToken, nil
}
