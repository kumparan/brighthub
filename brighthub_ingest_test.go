package brighthub

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icrowley/fake"

	"github.com/stretchr/testify/assert"
)

func TestClient_IngestVideo(t *testing.T) {
	httpMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"id": "id-video-lucu"}`)
	}))
	defer httpMock.Close()
	dynamicIngestBaseURL = httpMock.URL // change for test

	bh := newClientMock()
	bh.httpClient = httpMock.Client()

	resp, err := bh.IngestVideo("id-video-lucu", &IngestVideoRequest{
		Master:        &IngestVideoMaster{URL: fake.DomainName()},
		Priority:      PriorityNormal,
		CaptureImages: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "id-video-lucu", resp.ID)
}

func TestClient_GetIngestProfile(t *testing.T) {
	httpMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{
			"name": "multi-platform-standard-static",
			"display_name": "Multiplatform Standard",
			"description": "Deliver a wide range of content types across a variety of platforms on mobile and desktop.",
			"dynamic_origin": {
				"renditions": [
					"default/audio64",
					"default/audio128",
					"default/video700",
					"default/video2000",
					"default/video1700",
					"default/video1200",
					"default/audio96",
					"default/video450",
					"default/video900"
					]
			}
		}`)
	}))
	defer httpMock.Close()
	dynamicIngestBaseURL = httpMock.URL // change for test

	bh := newClientMock()
	bh.httpClient = httpMock.Client()

	resp, err := bh.GetIngestProfile("id-ingest-profile")
	assert.NoError(t, err)
	assert.Equal(t, "multi-platform-standard-static", resp.Name)
	assert.Equal(t, "Multiplatform Standard", resp.DisplayName)
	assert.Equal(t, "Deliver a wide range of content types across a variety of platforms on mobile and desktop.", resp.Description)
	assert.Equal(t, 9, len(resp.DynamicOrigin.Renditions))
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/audio64")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/audio128")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/video700")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/video2000")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/video1700")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/video1200")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/audio96")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/video450")
	assert.Contains(t, resp.DynamicOrigin.Renditions, "default/video900")
}
