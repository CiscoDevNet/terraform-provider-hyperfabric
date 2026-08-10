// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	frameworkdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSensitiveAttributesAreMarkedSensitive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	providerResponse := &frameworkprovider.SchemaResponse{}
	(&HyperfabricProvider{}).Schema(ctx, frameworkprovider.SchemaRequest{}, providerResponse)

	bearerTokenResourceResponse := &frameworkresource.SchemaResponse{}
	(&BearerTokenResource{}).Schema(ctx, frameworkresource.SchemaRequest{}, bearerTokenResourceResponse)

	resourceResponse := &frameworkresource.SchemaResponse{}
	(&NodeManagementPortResource{}).Schema(ctx, frameworkresource.SchemaRequest{}, resourceResponse)

	bearerTokenDataSourceResponse := &frameworkdatasource.SchemaResponse{}
	(&BearerTokenDataSource{}).Schema(ctx, frameworkdatasource.SchemaRequest{}, bearerTokenDataSourceResponse)

	dataSourceResponse := &frameworkdatasource.SchemaResponse{}
	(&NodeManagementPortDataSource{}).Schema(ctx, frameworkdatasource.SchemaRequest{}, dataSourceResponse)

	tests := map[string]bool{
		"provider API token":                              providerResponse.Schema.Attributes["token"].IsSensitive(),
		"provider proxy credentials":                      providerResponse.Schema.Attributes["proxy_creds"].IsSensitive(),
		"bearer token resource token":                     bearerTokenResourceResponse.Schema.Attributes["token"].IsSensitive(),
		"node management port proxy password":             resourceResponse.Schema.Attributes["proxy_password"].IsSensitive(),
		"bearer token data source token":                  bearerTokenDataSourceResponse.Schema.Attributes["token"].IsSensitive(),
		"node management port data source proxy password": dataSourceResponse.Schema.Attributes["proxy_password"].IsSensitive(),
	}

	for name, sensitive := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if !sensitive {
				t.Error("credential attribute must be marked sensitive")
			}
		})
	}
}
