package brighthub

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/icrowley/fake"

	"github.com/stretchr/testify/assert"
)

func newClientMock() *client {
	return &client{
		accountID:             fake.Characters(),
		clientID:              fake.Characters(),
		clientSecret:          fake.Characters(),
		accessToken:           fake.CharactersN(20),
		accessTokenAcquiredAt: time.Now(),
		httpClient:            defaultHTTPClient,
	}
}

func TestNew(t *testing.T) {
	httpMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{
		    "access_token": "kucing-lucu",
		    "token_type": "Bearer",
		    "expires_in": 300
		}`)
	}))
	defer httpMock.Close()
	authBaseURL = httpMock.URL // change for test

	bh, err := New("client-id", "client-secret", "account-id", httpMock.Client())
	assert.NoError(t, err)
	assert.NotNil(t, bh)

	bhc := bh.(*client)
	assert.Equal(t, "client-id", bhc.clientID)
	assert.Equal(t, "client-secret", bhc.clientSecret)
	assert.Equal(t, "account-id", bhc.accountID)
	assert.False(t, bhc.accessTokenAcquiredAt.IsZero())
	assert.NotEmpty(t, bhc.accessToken)
	assert.Equal(t, httpMock.Client(), bhc.httpClient)
}

func TestClient_getAccessToken(t *testing.T) {
	bhc := newClientMock()
	t.Run("token still valid", func(t *testing.T) {
		bhc.accessTokenAcquiredAt = time.Now()
		newToken, err := bhc.getAccessToken()
		assert.NoError(t, err)
		assert.Equal(t, bhc.accessToken, newToken)
	})

	t.Run("token already expired", func(t *testing.T) {
		httpMock2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{
		    "access_token": "kucing-oren",
		    "token_type": "Bearer",
		    "expires_in": 300
		}`)
		}))
		defer httpMock2.Close()
		authBaseURL = httpMock2.URL // change for test

		bhc.accessTokenAcquiredAt = time.Now().Add(-60 * time.Minute)
		bhc.httpClient = httpMock2.Client()
		newToken, err := bhc.getAccessToken()
		assert.NoError(t, err)
		assert.Equal(t, "kucing-oren", bhc.accessToken)
		assert.Equal(t, bhc.accessToken, newToken)
		assert.True(t, time.Since(bhc.accessTokenAcquiredAt).Minutes() < 1)
	})
}
