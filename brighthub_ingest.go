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
	// Priority ingest priority
	Priority string

	// IngestVideoRequest :nodoc:
	IngestVideoRequest struct {
		Master        *IngestVideoMaster `json:"master"`
		Priority      Priority           `json:"priority"`
		CaptureImages bool               `json:"capture-images"`
		Callbacks     []string           `json:"callbacks,omitempty"`
		// TODO add more request body
	}

	// IngestVideoMaster :nodoc:
	IngestVideoMaster struct {
		URL string `json:"url"`
	}

	// IngestVideoResponse :nodoc:
	IngestVideoResponse struct {
		ID string `json:"id"`
		// TODO add more response body
	}

	// IngestProfile :nodoc:
	IngestProfile struct {
		Name          string        `json:"name"`
		DisplayName   string        `json:"display_name"`
		Description   string        `json:"description"`
		DynamicOrigin DynamicOrigin `json:"dynamic_origin"`
	}

	// DynamicOrigin :nodoc:
	DynamicOrigin struct {
		Renditions []string `json:"renditions"`
	}
)

const (
	// PriorityLow :nodoc:
	PriorityLow Priority = "low"
	// PriorityNormal :nodoc:
	PriorityNormal Priority = "normal"
)

var (
	dynamicIngestBaseURL = "https://ingest.api.brightcove.com/v1"
	ingestionBaseURL     = "https://ingestion.api.brightcove.com/v1"
)

// IngestVideo :nodoc:
func (c *client) IngestVideo(videoID string, req *IngestVideoRequest) (*IngestVideoResponse, error) {
	token, err := c.getAccessToken()
	if err != nil {
		log.WithFields(log.Fields{
			"videoID": videoID,
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(&req)
	if err != nil {
		log.WithFields(log.Fields{
			"videoID": videoID,
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}
	r, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts/%s/videos/%s/ingest-requests", dynamicIngestBaseURL, c.accountID, videoID), b)
	if err != nil {
		log.WithFields(log.Fields{
			"videoID": videoID,
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		log.WithFields(log.Fields{
			"videoID": videoID,
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, ErrBadRequest
		case http.StatusUnauthorized:
			return nil, ErrUnauthorized
		case http.StatusForbidden:
			return nil, ErrDynamicDeliveryNotAllowed
		case http.StatusUnprocessableEntity:
			return nil, ErrIllegalField
		case http.StatusInternalServerError:
			return nil, ErrInternalError
		case http.StatusTooManyRequests:
			return nil, ErrRateLimitExceeded
		default:
			return nil, fmt.Errorf("undefined error with code %d", resp.StatusCode)
		}
	}

	ingestResponse := new(IngestVideoResponse)
	err = json.NewDecoder(resp.Body).Decode(&ingestResponse)
	if err != nil {
		log.WithFields(log.Fields{
			"videoID": videoID,
			"request": utils.Dump(req)}).
			Error(err)
		return nil, err
	}
	return ingestResponse, nil
}

// GetIngestProfile :nodoc:
func (c *client) GetIngestProfile(id string) (*IngestProfile, error) {
	token, err := c.getAccessToken()
	if err != nil {
		log.WithFields(log.Fields{"profileID": id}).Error(err)
		return nil, err
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/profiles/%s", ingestionBaseURL, c.accountID, id), nil)
	if err != nil {
		log.WithFields(log.Fields{"profileID": id}).Error(err)
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		log.WithFields(log.Fields{"profileID": id}).Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, ErrUnauthorized
		case http.StatusNotFound:
			return nil, ErrResourceNotFound
		case http.StatusConflict:
			return nil, ErrProfileError
		case http.StatusInternalServerError:
			return nil, ErrInternalError
		case http.StatusTooManyRequests:
			return nil, ErrRateLimitExceeded
		default:
			return nil, fmt.Errorf("undefined error with code %d", resp.StatusCode)
		}
	}

	ingestProfile := new(IngestProfile)
	err = json.NewDecoder(resp.Body).Decode(&ingestProfile)
	if err != nil {
		log.WithFields(log.Fields{"profileID": id}).Error(err)
		return nil, err
	}

	return ingestProfile, nil
}
