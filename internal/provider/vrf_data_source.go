// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &VrfDataSource{}

func NewVrfDataSource() datasource.DataSource {
	return &VrfDataSource{}
}

// VrfDataSource defines the data source implementation.
type VrfDataSource struct {
	client *client.Client
}

func (r *VrfDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_vrf")
	resp.TypeName = req.ProviderTypeName + "_vrf"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_vrf")
}

func (r *VrfDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_vrf")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "VRF data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a VRF in a Fabric.",
			},
			"vrf_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`vrf_id` defines the unique identifier of a VRF.",
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "`fabric_id` defines the unique identifier of a Fabric.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the VRF.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the VRF.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the VRF.",
				Computed:            true,
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "The flag that denote if the VRF is the default VRF or not.",
				Computed:            true,
			},
			"asn": schema.Float64Attribute{
				MarkdownDescription: "The Autonomous System Number (ASN) used for the VRF external connections.",
				Computed:            true,
			},
			"vni": schema.Float64Attribute{
				MarkdownDescription: "The VXLAN Network Identifier (VNI) used for the VRF.",
				Computed:            true,
			},
			"route_target": schema.StringAttribute{
				MarkdownDescription: "The route target associated with the VRF.",
				Computed:            true,
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsDataSourceSchemaAttribute(),
			"annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_vrf")
}

func (r *VrfDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_vrf")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_vrf")
}

func (r *VrfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_vrf")
	var data *VrfResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetNodeAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.VrfId = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_vrf with id '%s'", data.Id.ValueString()))

	getAndSetVrfAttributes(ctx, &resp.Diagnostics, r.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_node data source",
			fmt.Sprintf("The hyperfabric_node data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_vrf with id '%s'", data.Id.ValueString()))
}
