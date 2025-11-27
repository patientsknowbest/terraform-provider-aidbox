package aidbox

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestResponse struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestApiClient(t *testing.T) {
	t.Run("Should call the correct endpoint and deserialise to an interface", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("{\"name\": \"Lord Vetinari\", \"value\": 8}"))
		}))

		client := NewApiClient(server.URL, "foo", "bar")
		response := &TestResponse{}
		err := client.post(context.TODO(), "", "/endpoint", response)

		assert.Equal(t, nil, err)
		assert.Equal(t, TestResponse{"Lord Vetinari", 8}, *response)
	})

	t.Run("should log an error response that isn't JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(502)
			w.Write([]byte("Service Unavailable"))
		}))

		client := NewApiClient(server.URL, "foo", "bar")
		err := client.post(context.TODO(), "", "/endpoint", "")

		expectedError := fmt.Sprintf(`unexpected status code (502) received: 502 Bad Gateway

===== POST %s/endpoint =====

===== RESPONSE BODY =====
Service Unavailable
`, server.URL)

		assert.Equal(t, expectedError, err.Error())
	})

	t.Run("should pretty-print a JSON error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(422)
			w.Write([]byte("{ \"name\": \"Ankh-Morpork City Watch\", \"type\": \"Organisation\", \"members\": [\"Sam Vimes\", \"Fred Colon\", \"Nobby Nobbs\", \"Cheery Littlebottom\", \"Detritus\", \"Reg Shoe\"] }"))
		}))

		client := NewApiClient(server.URL, "foo", "bar")
		err := client.post(context.TODO(), "", "/endpoint", "")

		expectedError := fmt.Sprintf(`unexpected status code (422) received: 422 Unprocessable Entity

===== POST %s/endpoint =====

===== RESPONSE BODY =====
{
  "name": "Ankh-Morpork City Watch",
  "type": "Organisation",
  "members": [
    "Sam Vimes",
    "Fred Colon",
    "Nobby Nobbs",
    "Cheery Littlebottom",
    "Detritus",
    "Reg Shoe"
  ]
}
`, server.URL)

		assert.Equal(t, expectedError, err.Error())
	})
}
