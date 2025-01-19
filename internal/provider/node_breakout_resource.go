// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/cisco-open/terraform-provider-hyperfabric/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NodeBreakoutResource{}
var _ resource.ResourceWithImportState = &NodeBreakoutResource{}

func NewNodeBreakoutResource() resource.Resource {
	return &NodeBreakoutResource{}
}

// NodeBreakoutResource defines the resource implementation.
type NodeBreakoutResource struct {
	client *client.Client
}

// NodeBreakoutResourceModel describes the resource data model.
type NodeBreakoutResourceModel struct {
	Id          types.String `tfsdk:"id"`
	BreakoutId  types.String `tfsdk:"breakout_id"`
	NodeId      types.String `tfsdk:"node_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Breakouts   types.Set    `tfsdk:"breakouts"`
	Ports       types.Set    `tfsdk:"ports"`
	Mode        types.String `tfsdk:"mode"`
	Pluggable   types.String `tfsdk:"pluggable"`
	Metadata    types.Object `tfsdk:"metadata"`
	Labels      types.Set    `tfsdk:"labels"`
	Annotations types.Set    `tfsdk:"annotations"`
}

func getEmptyNodeBreakoutResourceModel() *NodeBreakoutResourceModel {
	return &NodeBreakoutResourceModel{
		Id:          basetypes.NewStringNull(),
		BreakoutId:  basetypes.NewStringNull(),
		NodeId:      basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		Enabled:     basetypes.NewBoolValue(false),
		Breakouts:   basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Ports:       basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Mode:        basetypes.NewStringNull(),
		Pluggable:   basetypes.NewStringNull(),
		Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewNodeBreakoutResourceModelFromData(data *NodeBreakoutResourceModel) *NodeBreakoutResourceModel {
	newNodeBreakout := getEmptyNodeBreakoutResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newNodeBreakout.Id = data.Id
	}

	if !data.BreakoutId.IsNull() && !data.BreakoutId.IsUnknown() {
		newNodeBreakout.BreakoutId = data.BreakoutId
	}

	if !data.NodeId.IsNull() && !data.NodeId.IsUnknown() {
		newNodeBreakout.NodeId = data.NodeId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newNodeBreakout.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newNodeBreakout.Description = data.Description
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		newNodeBreakout.Enabled = data.Enabled
	}

	if !data.Breakouts.IsNull() && !data.Breakouts.IsUnknown() {
		newNodeBreakout.Breakouts = data.Breakouts
	}

	if !data.Ports.IsNull() && !data.Ports.IsUnknown() {
		newNodeBreakout.Ports = data.Ports
	}

	if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
		newNodeBreakout.Mode = data.Mode
	}

	if !data.Pluggable.IsNull() && !data.Pluggable.IsUnknown() {
		newNodeBreakout.Pluggable = data.Pluggable
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newNodeBreakout.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newNodeBreakout.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newNodeBreakout.Annotations = data.Annotations
	}

	return newNodeBreakout
}

type NodeBreakoutIdentifier struct {
	Id types.String
}

func (r *NodeBreakoutResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_node_breakout")
	resp.TypeName = req.ProviderTypeName + "_node_breakout"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_node_breakout")
}

func (r *NodeBreakoutResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_node_breakout")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Breakouts resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Breakouts of a Node in a Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"breakout_id": schema.StringAttribute{
				MarkdownDescription: "`breakout_id` defines the unique identifier of a Breakouts of a Node in a Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"node_id": schema.StringAttribute{
				MarkdownDescription: "`node_id` defines the unique identifier of a Node in a Fabric.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Breakouts of the Node. The name should be in the `<Port Name>.<Integer>` format (i.e. `Ethernet1_1.100`). If `vlan_id` attribute is not provided, the integer in the Breakouts name will be used as the encapsulation VLAN ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Breakouts of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled admin state of the Breakouts of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"breakouts": getBreakoutBreakoutsSchemaAttribute(),
			"ports":     getBreakoutPortsSchemaAttribute(),
			"mode": schema.StringAttribute{
				MarkdownDescription: "The mode used to Breakout the Ports.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"pluggable": schema.StringAttribute{
				MarkdownDescription: "The type of pluggable used for the Breakout.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsSchemaAttribute(),
			"annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_node_breakout")
}

func getBreakoutBreakoutsSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of Breakout Port names of the Breakout of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getBreakoutBreakoutsDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of Breakout Port names of the Breakout of the Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func getBreakoutPortsSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of Node Ports names to be configured as Breakout Ports.`,
		Required:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		ElementType: types.StringType,
	}
}

func getBreakoutPortsDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of Node Ports names configured as Breakout Ports.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func (r *NodeBreakoutResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		var planData, stateData, configData *NodeBreakoutResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)

		if resp.Diagnostics.HasError() {
			return
		}

		if stateData != nil {
			// Set read-only fields in planData from stateData
			planData.Breakouts = stateData.Breakouts

			// Compare the string representation of the planData and stateData because structs cannot be compacted directly
			if fmt.Sprintf("%s", planData) != fmt.Sprintf("%s", stateData) {
				planData.Breakouts = basetypes.NewSetUnknown(types.StringType)
			}
		}

		resp.Diagnostics.Append(resp.Plan.Set(ctx, &planData)...)
	}
}

func (r *NodeBreakoutResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_node_breakout")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
	tflog.Debug(ctx, "End configure of resource: hyperfabric_node_breakout")
}

func (r *NodeBreakoutResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_node_breakout")

	var data *NodeBreakoutResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_node_breakout with name '%s'", data.Name.ValueString()))

	jsonPayload := getNodeBreakoutJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	container := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/breakouts", data.NodeId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	breakoutsContainer, err := container.ArrayElement(0, "breakouts")
	if err != nil {
		return
	}

	breakoutsId := StripQuotes(breakoutsContainer.Search("id").String())
	if breakoutsId != "" {
		data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/breakouts/%s", data.NodeId.ValueString(), breakoutsId))
		data.BreakoutId = basetypes.NewStringValue(breakoutsId)
		getAndSetNodeBreakoutAttributes(ctx, &resp.Diagnostics, r.client, data)
	} else {
		data.Id = basetypes.NewStringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))
}

func (r *NodeBreakoutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_node_breakout")
	var data *NodeBreakoutResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))
	checkAndSetNodeBreakoutIds(data)
	getAndSetNodeBreakoutAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	if data.Id.IsNull() {
		var emptyData *NodeBreakoutResourceModel
		resp.Diagnostics.Append(resp.State.Set(ctx, &emptyData)...)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))
}

func (r *NodeBreakoutResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_node_breakout")
	var data *NodeBreakoutResourceModel
	var stateData *NodeBreakoutResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))

	jsonPayload := getNodeBreakoutJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/breakouts/%s", data.NodeId.ValueString(), data.BreakoutId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetNodeBreakoutAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))
}

func (r *NodeBreakoutResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_node_breakout")
	var data *NodeBreakoutResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))
	checkAndSetNodeBreakoutIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/breakouts/%s", data.NodeId.ValueString(), data.BreakoutId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_node_breakout with id '%s'", data.Id.ValueString()))
}

func (r *NodeBreakoutResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_node_breakout")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *NodeBreakoutResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_node_breakout with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_node_breakout")
}

func getAndSetNodeBreakoutAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *NodeBreakoutResourceModel) {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/breakouts/%s", data.NodeId.ValueString(), data.BreakoutId.ValueString()), "GET", nil)
	if diags.HasError() {
		return
	}

	newNodeBreakout := *getNewNodeBreakoutResourceModelFromData(data)
	node := getEmptyNodeResourceModel()
	node.Id = newNodeBreakout.NodeId
	checkAndSetNodeIds(node)

	if requestData.Data() != nil {
		for attributeName, attributeValue := range requestData.Data().(map[string]interface{}) {
			if attributeName == "id" && (data.BreakoutId.IsNull() || data.BreakoutId.IsUnknown() || data.BreakoutId.ValueString() == "" || data.BreakoutId.ValueString() != attributeValue.(string)) {
				newNodeBreakout.BreakoutId = basetypes.NewStringValue(attributeValue.(string))
				newNodeBreakout.Id = basetypes.NewStringValue(fmt.Sprintf("%s/breakouts/%s", newNodeBreakout.NodeId.ValueString(), newNodeBreakout.BreakoutId.ValueString()))
			} else if attributeName == "fabricId" && (node.FabricId.IsNull() || node.FabricId.IsUnknown() || node.FabricId.ValueString() == "" || node.FabricId.ValueString() != attributeValue.(string)) {
				node.FabricId = basetypes.NewStringValue(attributeValue.(string))
				newNodeBreakout.NodeId = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", node.FabricId.ValueString(), node.NodeId.ValueString()))
				newNodeBreakout.Id = basetypes.NewStringValue(fmt.Sprintf("%s/breakouts/%s", newNodeBreakout.NodeId.ValueString(), newNodeBreakout.BreakoutId.ValueString()))
			} else if attributeName == "nodeId" && (node.NodeId.IsNull() || node.NodeId.IsUnknown() || node.NodeId.ValueString() == "" || node.NodeId.ValueString() != attributeValue.(string)) {
				node.NodeId = basetypes.NewStringValue(attributeValue.(string))
				newNodeBreakout.NodeId = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", node.FabricId.ValueString(), node.NodeId.ValueString()))
				newNodeBreakout.Id = basetypes.NewStringValue(fmt.Sprintf("%s/breakouts/%s", newNodeBreakout.NodeId.ValueString(), newNodeBreakout.BreakoutId.ValueString()))
			} else if attributeName == "name" {
				newNodeBreakout.Name = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "description" {
				newNodeBreakout.Description = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "enabled" {
				newNodeBreakout.Enabled = basetypes.NewBoolValue(attributeValue.(bool))
			} else if attributeName == "breakouts" {
				newNodeBreakout.Breakouts = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "ports" {
				newNodeBreakout.Ports = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "mode" {
				newNodeBreakout.Mode = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "pluggable" {
				newNodeBreakout.Pluggable = basetypes.NewStringValue(attributeValue.(string))
			} else if attributeName == "metadata" {
				newNodeBreakout.Metadata = NewMetadataObject(ctx, attributeValue.(map[string]interface{}))
			} else if attributeName == "labels" {
				newNodeBreakout.Labels = NewSetString(ctx, attributeValue.([]interface{}))
			} else if attributeName == "annotations" {
				newNodeBreakout.Annotations = NewAnnotationsSet(ctx, attributeValue.([]interface{}))
			}
		}
	} else {
		newNodeBreakout.Id = basetypes.NewStringNull()
	}
	*data = newNodeBreakout
}

func getNodeBreakoutJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *NodeBreakoutResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	payloadMap["enabled"] = true

	if !data.Ports.IsNull() && !data.Ports.IsUnknown() {
		payloadMap["ports"] = getSetStringJsonPayload(ctx, data.Ports)
	}

	if !data.Mode.IsNull() && !data.Mode.IsUnknown() {
		payloadMap["mode"] = data.Mode.ValueString()
	}

	if !data.Pluggable.IsNull() && !data.Pluggable.IsUnknown() {
		payloadMap["pluggable"] = data.Pluggable.ValueString()
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		payloadMap["labels"] = getSetStringJsonPayload(ctx, data.Labels)
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		payloadMap["annotations"] = getAnnotationsJsonPayload(ctx, data.Annotations)
	}

	var payload map[string]interface{}
	if action == "create" {
		payloadList = append(payloadList, payloadMap)
		payload = map[string]interface{}{"breakouts": payloadList}
	} else {
		payload = payloadMap
	}

	marshalPayload, err := json.Marshal(payload)
	if err != nil {
		diags.AddError(
			"Marshalling of JSON payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}

	jsonPayload, err := gabs.ParseJSON(marshalPayload)
	if err != nil {
		diags.AddError(
			"Construction of JSON payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}
	return jsonPayload
}

func checkAndSetNodeBreakoutIds(data *NodeBreakoutResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/breakouts/") {
		if data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" ||
			data.BreakoutId.IsNull() || data.BreakoutId.IsUnknown() || data.BreakoutId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/breakouts/")
			data.NodeId = basetypes.NewStringValue(splitId[0])
			data.BreakoutId = basetypes.NewStringValue(splitId[1])
		}
	}
}
