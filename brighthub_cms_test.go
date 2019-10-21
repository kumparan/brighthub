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
		w.WriteHeader(http.StatusCreated)
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

func TestClient_GetVideoMasterInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{
				"encoding_rate": 23152000,
  				"height": 1080,
  				"width": 1920,
  				"id": "a0a2e032-4de4-4495-a59e-a806d52989",
  				"size": 90990884,
  				"updated_at": "2019-04-30T10:09:12.548Z",
  				"created_at": "2019-04-30T10:09:12.548Z",
  				"duration": 31431
			}`)
		}))
		defer httpMock.Close()
		cmsBaseURL = httpMock.URL // change for test

		bh := newClientMock()
		bh.httpClient = httpMock.Client()

		videoMasterInfo, err := bh.GetVideoMasterInfo("12345")
		assert.NoError(t, err)
		assert.EqualValues(t, int64(23152000), videoMasterInfo.EncodingRate)
		assert.EqualValues(t, int64(1080), videoMasterInfo.Height)
		assert.EqualValues(t, int64(1920), videoMasterInfo.Width)
		assert.EqualValues(t, "a0a2e032-4de4-4495-a59e-a806d52989", videoMasterInfo.ID)
		assert.EqualValues(t, int64(90990884), videoMasterInfo.Size)
		assert.EqualValues(t, "2019-04-30T10:09:12.548Z", videoMasterInfo.CreatedAt)
		assert.EqualValues(t, "2019-04-30T10:09:12.548Z", videoMasterInfo.UpdatedAt)
		assert.EqualValues(t, int64(31431), videoMasterInfo.Duration)
	})
}
