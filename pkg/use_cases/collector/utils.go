package collector

import (
	"net/http"
	"time"

	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
)

// Creates the HTTP header for requesting a token.
func MakesTokenRequestHeader(setup entities.CopergasConfig) http.Header {
	header := http.Header{
		"Accept":     []string{"application/json"},
		"grant_type": []string{"password"},
		"username":   []string{setup.Credentials.Username},
		"password":   []string{setup.Credentials.Password},
	}

	return header
}

func Wait(seconds float32) {
	timeToWaitInSeconds := seconds
	time.Sleep(time.Duration(timeToWaitInSeconds) * time.Second)
}

func CreatesHTTPClient(seconds float32) *http.Client {
	httpClientTimeoutInMinutes := int((seconds / 60.0) / 2.0)
	return &http.Client{Timeout: time.Duration(httpClientTimeoutInMinutes) * time.Minute}
}
