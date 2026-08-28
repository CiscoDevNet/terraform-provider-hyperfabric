// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoRestRequestNotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusNotFound)
		_, _ = response.Write([]byte(`{"errCode":"NOT_FOUND","message":"missing"}`))
	}))
	t.Cleanup(server.Close)

	restClient := NewClient(server.URL, "test-token", HttpClient(server.Client()), SkipLoggingPayload(true))

	for _, method := range []string{http.MethodGet, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			result, diagnostic := restClient.DoRestRequest("/objects/missing", method, nil)

			if diagnostic != nil {
				t.Fatalf("expected no diagnostic, got: %s: %s", diagnostic.Summary, diagnostic.Detail)
			}
			if result == nil {
				t.Fatal("expected a non-nil result")
			}
			if result.Found {
				t.Fatal("expected object to be reported as not found")
			}
			if result.StatusCode != http.StatusNotFound {
				t.Fatalf("expected status %d, got %d", http.StatusNotFound, result.StatusCode)
			}
		})
	}
}

func TestDoRestRequestDecode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"id":"object-1"}`))
	}))
	t.Cleanup(server.Close)

	restClient := NewClient(server.URL, "test-token", HttpClient(server.Client()), SkipLoggingPayload(true))
	result, diagnostic := restClient.DoRestRequest("/objects/object-1", http.MethodGet, nil)
	if diagnostic != nil {
		t.Fatalf("expected no diagnostic, got: %s: %s", diagnostic.Summary, diagnostic.Detail)
	}

	var response struct {
		ID string `json:"id"`
	}
	if err := result.Decode(&response); err != nil {
		t.Fatalf("expected response to decode: %s", err)
	}
	if response.ID != "object-1" {
		t.Fatalf("expected object-1, got %q", response.ID)
	}
}

func TestDoRestRequestTransportError(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection failed")
	})}
	restClient := NewClient("https://hyperfabric.example.com", "test-token", HttpClient(httpClient), SkipLoggingPayload(true))

	result, diagnostic := restClient.DoRestRequest("/objects/object-1", http.MethodGet, nil)
	if result != nil {
		t.Fatalf("expected no result, got %#v", result)
	}
	if diagnostic == nil {
		t.Fatal("expected a diagnostic")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
