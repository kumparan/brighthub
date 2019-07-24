package brighthub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kumparan/go-lib/utils"
	log "github.com/sirupsen/logrus"
)

type (
	// State :nodoc:
	State string

	// CreateVideoRequest :nodoc:
	CreateVideoRequest struct {
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		LongDescription string   `json:"long_description"`
		ReferenceID     string   `json:"reference_id,omitempty"`
		State           State    `json:"state"`
		Tags            []string `json:"tags,omitempty"`
		// TODO Add more request body
		// to richest create video request
	}

	// CreateVideoResponse :nodoc:
	CreateVideoResponse struct {
		ID        string `json:"id"`
		AccountID string `json:"account_id"`
		// TODO add more response field
	}
)

const (
	// StateActive :nodoc:
	StateActive State = "ACTIVE"
	// StateInactive :nodoc:
	StateInactive State = "INACTIVE"
)

var cmsBaseURL = "https://cms.api.brightcove.com/v1"

// CreateVideo :nodoc:
func (c *client) CreateVideo(req *CreateVideoRequest) (*CreateVideoResponse, error) {
	token, err := c.getAccessToken()
	if err != nil {
		log.WithFields(log.Fields{
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(&req)
	if err != nil {
		log.WithFields(log.Fields{
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}
	r, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts/%s/videos", cmsBaseURL, c.accountID), b)
	if err != nil {
		log.WithFields(log.Fields{
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		log.WithFields(log.Fields{
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, ErrUnauthorized
		case http.StatusForbidden, http.StatusUnprocessableEntity:
			return nil, ErrIllegalField
		case http.StatusMethodNotAllowed:
			return nil, ErrMethodNotAllowed
		case http.StatusConflict:
			return nil, ErrDuplicateReferenceID
		case http.StatusTooManyRequests:
			return nil, ErrTooManyRequest
		default:
			return nil, fmt.Errorf("undefined error with code %d", resp.StatusCode)
		}
	}

	videoResponse := new(CreateVideoResponse)
	err = json.NewDecoder(resp.Body).Decode(&videoResponse)
	if err != nil {
		log.WithFields(log.Fields{
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}

	return videoResponse, nil
}

// AddVideoToFolder :nodoc:
func (c *client) AddVideoToFolder(videoID, folderID string) error {
	token, err := c.getAccessToken()
	if err != nil {
		log.WithFields(log.Fields{
			"folderID": folderID,
			"videoID":  videoID}).
			Error(err)
		return err
	}

	r, err := http.NewRequest("PUT", fmt.Sprintf("%s/accounts/%s/folders/%s/videos/%s", cmsBaseURL, c.accountID, folderID, videoID), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"folderID": folderID,
			"request":  videoID}).
			Error(err)
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		log.WithFields(log.Fields{
			"folderID": folderID,
			"videoID":  videoID}).
			Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusForbidden:
			return ErrNotAvailable
		case http.StatusNotFound:
			return ErrResourceNotFound
		case http.StatusMethodNotAllowed:
			return ErrMethodNotAllowed
		case http.StatusTooManyRequests:
			return ErrTooManyRequest
		case http.StatusInternalServerError:
			return ErrInternalError
		default:
			return fmt.Errorf("undefined error with code %d", resp.StatusCode)
		}
	}

	return nil
}
