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

	"github.com/CiscoDevNet/terraform-provider-hyperfabric/internal/client"
	customTypes "github.com/CiscoDevNet/terraform-provider-hyperfabric/internal/provider/custom_types"
	"github.com/Jeffail/gabs/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NodeLoopbackResource{}
var _ resource.ResourceWithImportState = &NodeLoopbackResource{}

func NewNodeLoopbackResource() resource.Resource {
	return &NodeLoopbackResource{}
}

// NodeLoopbackResource defines the resource implementation.
type NodeLoopbackResource struct {
	client *client.Client
}

// NodeLoopbackResourceModel describes the resource data model.
type NodeLoopbackResourceModel struct {
	Id          types.String `tfsdk:"id"`
	NodeId      types.String `tfsdk:"node_id"`
	LoopbackId  types.String `tfsdk:"loopback_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	// Enabled     types.Bool   `tfsdk:"enabled"`
	Ipv4Address types.String                      `tfsdk:"ipv4_address"`
	Ipv6Address types.String                      `tfsdk:"ipv6_address"`
	VrfId       customTypes.UuidFromIdStringValue `tfsdk:"vrf_id"`
	Metadata    types.Object                      `tfsdk:"metadata"`
	Labels      types.Set                         `tfsdk:"labels"`
	Annotations types.Set                         `tfsdk:"annotations"`
}

func getEmptyNodeLoopbackResourceModel() *NodeLoopbackResourceModel {
	return &NodeLoopbackResourceModel{
		Id:          basetypes.NewStringNull(),
		NodeId:      basetypes.NewStringNull(),
		LoopbackId:  basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		// Enabled:     basetypes.NewBoolValue(false),
		Ipv4Address: basetypes.NewStringNull(),
		Ipv6Address: basetypes.NewStringNull(),
		VrfId:       customTypes.NewUuidFromIdStringNull(),
		Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewNodeLoopbackResourceModelFromData(data *NodeLoopbackResourceModel) *NodeLoopbackResourceModel {
	newNodeLoopback := getEmptyNodeLoopbackResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newNodeLoopback.Id = data.Id
	}

	if !data.LoopbackId.IsNull() && !data.LoopbackId.IsUnknown() {
		newNodeLoopback.LoopbackId = data.LoopbackId
	}

	if !data.NodeId.IsNull() && !data.NodeId.IsUnknown() {
		newNodeLoopback.NodeId = data.NodeId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newNodeLoopback.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newNodeLoopback.Description = data.Description
	}

	// if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
	// 	newNodeLoopback.Enabled = data.Enabled
	// }

	if !data.Ipv4Address.IsNull() && !data.Ipv4Address.IsUnknown() {
		newNodeLoopback.Ipv4Address = data.Ipv4Address
	}

	if !data.Ipv6Address.IsNull() && !data.Ipv6Address.IsUnknown() {
		newNodeLoopback.Ipv6Address = data.Ipv6Address
	}

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		newNodeLoopback.VrfId = data.VrfId
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newNodeLoopback.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newNodeLoopback.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newNodeLoopback.Annotations = data.Annotations
	}

	return newNodeLoopback
}

type NodeLoopbackIdentifier struct {
	Id types.String
}

func (r *NodeLoopbackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_node_loopback")
	resp.TypeName = req.ProviderTypeName + "_node_loopback"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_node_loopback")
}

func (r *NodeLoopbackResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_node_loopback")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node Loopback resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "`id` defines the unique identifier of the Loopback of a Node in a Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"loopback_id": schema.StringAttribute{
				MarkdownDescription: "`loopback_id` defines the unique identifier of a Loopback of a Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
				MarkdownDescription: "The name of the Loopback of the Node.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Loopback of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			// "enabled": schema.BoolAttribute{
			// 	MarkdownDescription: "The enabled admin state of the Loopback of the Node.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	// Default:             booldefault.StaticBool(true),
			// 	PlanModifiers: []planmodifier.Bool{
			// 		boolplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			"ipv4_address": schema.StringAttribute{
				MarkdownDescription: "The IPv4 address configured on the Loopback of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"ipv6_address": schema.StringAttribute{
				MarkdownDescription: "The IPv6 address configured on the Loopback of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"vrf_id": schema.StringAttribute{
				CustomType:          customTypes.UuidFromIdStringType{},
				MarkdownDescription: "The `vrf_id` of a VRF to associate with the Loopback of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
					CompareUuidWithIdForEquality(),
				},
			},
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsSchemaAttribute(),
			"annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_node_loopback")
}

func (r *NodeLoopbackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_node_loopback")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_node_loopback")
}

func (r *NodeLoopbackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_node_loopback")

	var data *NodeLoopbackResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_node_loopback with name '%s'", data.Name.ValueString()))

	jsonPayload := getNodeLoopbackJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	result := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks", data.NodeId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	var response nodeLoopbacksAPIResponse
	if !decodeRestResult(&resp.Diagnostics, result, &response, "node loopback create") {
		return
	}
	createdLoopback, ok := requireFirstAPIObject(&resp.Diagnostics, response.Loopbacks, "created node loopback")
	if !ok {
		return
	}
	loopbackID, ok := requireAPIIdentifier(&resp.Diagnostics, createdLoopback.Id, "created node loopback")
	if !ok {
		return
	}

	data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/loopbacks/%s", data.NodeId.ValueString(), loopbackID))
	data.LoopbackId = basetypes.NewStringValue(loopbackID)
	getAndSetNodeLoopbackAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_node_loopback")
	var data *NodeLoopbackResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
	checkAndSetNodeLoopbackIds(data)
	found := getAndSetNodeLoopbackAttributes(ctx, &resp.Diagnostics, r.client, data)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_node_loopback")
	var data *NodeLoopbackResourceModel
	var stateData *NodeLoopbackResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))

	jsonPayload := getNodeLoopbackJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks/%s", data.NodeId.ValueString(), data.LoopbackId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetNodeLoopbackAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_node_loopback")
	var data *NodeLoopbackResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
	checkAndSetNodeLoopbackIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks/%s", data.NodeId.ValueString(), data.LoopbackId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_node_loopback with id '%s'", data.Id.ValueString()))
}

func (r *NodeLoopbackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_node_loopback")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *NodeLoopbackResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_node_loopback with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_node_loopback")
}

func getAndSetNodeLoopbackAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *NodeLoopbackResourceModel) bool {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/loopbacks/%s", data.NodeId.ValueString(), data.LoopbackId.ValueString()), "GET", nil)
	if diags.HasError() {
		return false
	}

	newNodeLoopback := *getNewNodeLoopbackResourceModelFromData(data)
	node := getEmptyNodeResourceModel()
	node.Id = newNodeLoopback.NodeId
	checkAndSetNodeIds(node)

	if requestData == nil || !requestData.Found {
		newNodeLoopback.Id = basetypes.NewStringNull()
		*data = newNodeLoopback
		return false
	}

	var response nodeLoopbackAPIResponse
	if !decodeRestResult(diags, requestData, &response, "node loopback") {
		return true
	}
	loopbackID, ok := requireAPIIdentifier(diags, response.Id, "node loopback")
	if !ok {
		return true
	}
	newNodeLoopback.LoopbackId = basetypes.NewStringValue(loopbackID)
	if response.FabricId != nil {
		node.FabricId = basetypes.NewStringValue(*response.FabricId)
	}
	if response.NodeId != nil {
		node.NodeId = basetypes.NewStringValue(*response.NodeId)
	}
	newNodeLoopback.NodeId = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", node.FabricId.ValueString(), node.NodeId.ValueString()))
	newNodeLoopback.Id = basetypes.NewStringValue(fmt.Sprintf("%s/loopbacks/%s", newNodeLoopback.NodeId.ValueString(), newNodeLoopback.LoopbackId.ValueString()))
	if response.Name != nil {
		newNodeLoopback.Name = basetypes.NewStringValue(*response.Name)
	}
	if response.Description != nil {
		newNodeLoopback.Description = basetypes.NewStringValue(*response.Description)
	}
	if response.IPv4Address != nil {
		newNodeLoopback.Ipv4Address = basetypes.NewStringValue(*response.IPv4Address)
	}
	if response.IPv6Address != nil {
		newNodeLoopback.Ipv6Address = basetypes.NewStringValue(*response.IPv6Address)
	}
	if response.VrfID != nil {
		newNodeLoopback.VrfId = customTypes.NewUuidFromIdStringValue(*response.VrfID)
	}
	if response.Metadata != nil {
		newNodeLoopback.Metadata = NewMetadataObject(ctx, response.Metadata)
	}
	if response.Labels != nil {
		newNodeLoopback.Labels = NewSetString(ctx, response.Labels)
	}
	if response.Annotations != nil {
		newNodeLoopback.Annotations = NewAnnotationsSet(ctx, response.Annotations)
	}
	*data = newNodeLoopback
	return true
}

func getNodeLoopbackJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *NodeLoopbackResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	// if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
	// 	payloadMap["enabled"] = data.Enabled.ValueBool()
	// }
	payloadMap["enabled"] = true

	if !data.Ipv4Address.IsNull() && !data.Ipv4Address.IsUnknown() {
		payloadMap["ipv4Address"] = data.Ipv4Address.ValueString()
	}

	if !data.Ipv6Address.IsNull() && !data.Ipv6Address.IsUnknown() {
		payloadMap["ipv6Address"] = data.Ipv6Address.ValueString()
	}

	if !data.VrfId.IsNull() && !data.VrfId.IsUnknown() {
		payloadMap["vrfId"] = data.VrfId.ValueString()
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
		payload = map[string]interface{}{"loopbacks": payloadList}
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

func checkAndSetNodeLoopbackIds(data *NodeLoopbackResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/loopbacks/") {
		if data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" ||
			data.LoopbackId.IsNull() || data.LoopbackId.IsUnknown() || data.LoopbackId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/loopbacks/")
			data.NodeId = basetypes.NewStringValue(splitId[0])
			data.LoopbackId = basetypes.NewStringValue(splitId[1])
		}
	}
}
