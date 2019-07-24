package brighthub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client interface {
	AddVideoToFolder(videoID, folderID string) error
	CreateVideo(req *CreateVideoRequest) (*CreateVideoResponse, error)
	IngestVideo(videoID string, req *IngestVideoRequest) (*IngestVideoResponse, error)
}

type client struct {
	accessToken  string
	accountID    string
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

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

const authBaseURL = "https://oauth.brightcove.com/v4"

type getAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *client) getAccessToken() (string, error) {
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
	return a.AccessToken, nil
}
