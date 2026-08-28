// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestGetAndSetFabricAttributesNotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusNotFound)
		_, _ = response.Write([]byte(`{"errCode":"NOT_FOUND","message":"missing"}`))
	}))
	t.Cleanup(server.Close)

	restClient := client.NewClient(server.URL, "test-token", client.HttpClient(server.Client()), client.SkipLoggingPayload(true))
	data := getEmptyFabricResourceModel()
	data.Id = basetypes.NewStringValue("fabric-1")
	var diagnostics diag.Diagnostics

	found := getAndSetFabricAttributes(context.Background(), &diagnostics, restClient, data)

	if found {
		t.Fatal("expected fabric to be reported as not found")
	}
	if diagnostics.HasError() {
		t.Fatalf("expected no diagnostic for a missing resource, got: %v", diagnostics)
	}
	if !data.Id.IsNull() {
		t.Fatalf("expected the resource identifier to be null, got: %s", data.Id)
	}
}

func TestGetAndSetFabricAttributesInvalidResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"fabricId":42}`))
	}))
	t.Cleanup(server.Close)

	restClient := client.NewClient(server.URL, "test-token", client.HttpClient(server.Client()), client.SkipLoggingPayload(true))
	data := getEmptyFabricResourceModel()
	data.Id = basetypes.NewStringValue("fabric-1")
	var diagnostics diag.Diagnostics

	found := getAndSetFabricAttributes(context.Background(), &diagnostics, restClient, data)

	if !found {
		t.Fatal("expected the API object to be reported as found")
	}
	if !diagnostics.HasError() {
		t.Fatal("expected an invalid API response diagnostic")
	}
}
