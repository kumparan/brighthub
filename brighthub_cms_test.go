package brighthub

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestClient_CreateVideo(t *testing.T) {
	httpMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"id": "id-video-lucu", "account_id": "account-id-kamu"}`)
	}))
	defer httpMock.Close()
	cmsBaseURL = httpMock.URL // change for test

	bh := newClientMock()
	bh.httpClient = httpMock.Client()

	resp, err := bh.CreateVideo(&CreateVideoRequest{
		Name:            fake.Title(),
		Description:     fake.Paragraph(),
		LongDescription: fake.Paragraphs(),
		ReferenceID:     fake.Characters(),
		State:           StateActive,
		Tags:            []string{fake.Model()},
	})
	assert.NoError(t, err)
	assert.Equal(t, "id-video-lucu", resp.ID)
	assert.Equal(t, "account-id-kamu", resp.AccountID)
}

func TestClient_AddVideoToFolder(t *testing.T) {
	httpMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer httpMock.Close()
	cmsBaseURL = httpMock.URL // change for test

	bh := newClientMock()
	bh.httpClient = httpMock.Client()

	err := bh.AddVideoToFolder("id-video-lucu", "id-folder-lucu")
	assert.NoError(t, err)
}
