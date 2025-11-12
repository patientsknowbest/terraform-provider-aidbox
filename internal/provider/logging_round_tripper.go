package provider

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// LoggingRoundTripper is an http.RoundTripper that logs the request and response
type LoggingRoundTripper struct {
	// Proxied is the RoundTripper to which the request is delegated.
	Proxied http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log request
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("Error dumping request: %v\n", err)
	} else {
		fmt.Printf("--- REQUEST ---\n%s\n", string(requestDump))
	}

	// Delegate the request to the original RoundTripper
	resp, err := lrt.Proxied.RoundTrip(req)

	if err != nil {
		fmt.Printf("--- RESPONSE ERROR ---\nError making request: %v\n", err)
		return nil, err
	}

	// Log response
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Printf("Error dumping response: %v\n", err)
	} else {
		fmt.Printf("--- RESPONSE ---\n%s\n", string(responseDump))
	}

	return resp, err
}
