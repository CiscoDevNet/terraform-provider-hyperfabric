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
var _ datasource.DataSource = &NodeBreakoutDataSource{}

func NewNodeBreakoutDataSource() datasource.DataSource {
	return &NodeBreakoutDataSource{}
}

// NodeBreakoutDataSource defines the data source implementation.
type NodeBreakoutDataSource struct {
	client *client.Client
}

func (d *NodeBreakoutDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of datasource: hyperfabric_node_breakout")
	resp.TypeName = req.ProviderTypeName + "_node_breakout"
	tflog.Debug(ctx, "End metadata of datasource: hyperfabric_node_breakout")
}

func (d *NodeBreakoutDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of datasource: hyperfabric_node_breakout")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Breakout data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Breakout of a Node in a Fabric.",
				Computed:            true,
			},
			"breakout_id": schema.StringAttribute{
				MarkdownDescription: "`breakout_id` defines the unique identifier of a Breakout of a Node in a Fabric.",
				Computed:            true,
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node in a Fabric.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Breakout of the Node.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Breakout of the Node.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled admin state of the Breakout of the Node.",
				Computed:            true,
			},
			"breakouts": getBreakoutBreakoutsDataSourceSchemaAttribute(),
			"ports":     getBreakoutPortsDataSourceSchemaAttribute(),
			"mode": schema.StringAttribute{
				MarkdownDescription: "The mode used to Breakout the Ports.",
				Computed:            true,
			},
			"pluggable": schema.StringAttribute{
				MarkdownDescription: "The type of pluggable used for the Breakout.",
				Computed:            true,
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsDataSourceSchemaAttribute(),
			"annotations": getAnnotationsDataSourceSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of datasource: hyperfabric_node_breakout")
}

func (d *NodeBreakoutDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of datasource: hyperfabric_node_breakout")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
	tflog.Debug(ctx, "End configure of datasource: hyperfabric_node_breakout")
}

func (d *NodeBreakoutDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Start read of datasource: hyperfabric_node_breakout")
	var data *NodeBreakoutResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a copy of the Id for when not found during getAndSetNodeBreakoutAttributes
	cachedId := data.Id.ValueString()
	if cachedId == "" && data.Name.ValueString() != "" {
		data.BreakoutId = data.Name
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of datasource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))

	getAndSetNodeBreakoutAttributes(ctx, &resp.Diagnostics, d.client, data)

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Failed to read hyperfabric_node_breakout data source",
			fmt.Sprintf("The hyperfabric_node_breakout data source with id '%s' has not been found", cachedId),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of datasource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))
}
