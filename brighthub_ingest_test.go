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
