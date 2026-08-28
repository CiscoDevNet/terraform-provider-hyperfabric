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
	"github.com/Jeffail/gabs/v2"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NodeResource{}
var _ resource.ResourceWithImportState = &NodeResource{}

func NewNodeResource() resource.Resource {
	return &NodeResource{}
}

// NodeResource defines the resource implementation.
type NodeResource struct {
	client *client.Client
}

// NodeResourceModel describes the resource data model.
type NodeResourceModel struct {
	Id           types.String `tfsdk:"id"`
	NodeId       types.String `tfsdk:"node_id"`
	FabricId     types.String `tfsdk:"fabric_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Location     types.String `tfsdk:"location"`
	ModelName    types.String `tfsdk:"model_name"`
	SerialNumber types.String `tfsdk:"serial_number"`
	DeviceId     types.String `tfsdk:"device_id"`
	Position     types.String `tfsdk:"position"`
	Roles        types.Set    `tfsdk:"roles"`
	Metadata     types.Object `tfsdk:"metadata"`
	Labels       types.Set    `tfsdk:"labels"`
	Annotations  types.Set    `tfsdk:"annotations"`
}

func getEmptyNodeResourceModel() *NodeResourceModel {
	return &NodeResourceModel{
		Id:           basetypes.NewStringNull(),
		NodeId:       basetypes.NewStringNull(),
		FabricId:     basetypes.NewStringNull(),
		Name:         basetypes.NewStringNull(),
		Description:  basetypes.NewStringNull(),
		Enabled:      basetypes.NewBoolValue(false),
		Location:     basetypes.NewStringNull(),
		ModelName:    basetypes.NewStringNull(),
		SerialNumber: basetypes.NewStringNull(),
		DeviceId:     basetypes.NewStringNull(),
		Position:     basetypes.NewStringNull(),
		Roles:        basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Metadata:     basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:       basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations:  basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewNodeResourceModelFromData(data *NodeResourceModel) *NodeResourceModel {
	newNode := getEmptyNodeResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newNode.Id = data.Id
	}

	if !data.NodeId.IsNull() && !data.NodeId.IsUnknown() {
		newNode.NodeId = data.NodeId
	}

	if !data.FabricId.IsNull() && !data.FabricId.IsUnknown() {
		newNode.FabricId = data.FabricId
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newNode.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newNode.Description = data.Description
	}

	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		newNode.Enabled = data.Enabled
	}

	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		newNode.Location = data.Location
	}

	if !data.ModelName.IsNull() && !data.ModelName.IsUnknown() {
		newNode.ModelName = data.ModelName
	}

	if !data.SerialNumber.IsNull() && !data.SerialNumber.IsUnknown() {
		newNode.SerialNumber = data.SerialNumber
	}

	if !data.DeviceId.IsNull() && !data.DeviceId.IsUnknown() {
		newNode.DeviceId = data.DeviceId
	}

	if !data.Position.IsNull() && !data.Position.IsUnknown() {
		newNode.Position = data.Position
	}

	if !data.Roles.IsNull() && !data.Roles.IsUnknown() {
		newNode.Roles = data.Roles
	}

	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		newNode.Metadata = data.Metadata
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newNode.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newNode.Annotations = data.Annotations
	}

	return newNode
}

type NodeIdentifier struct {
	Id types.String
}

func (r *NodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_node")
	resp.TypeName = req.ProviderTypeName + "_node"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_node")
}

func (r *NodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_node")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Node resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a Node in a Fabric.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`node_id` defines the unique identifier of a Node.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "`fabric_id` defines the unique identifier of a Fabric.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Node.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "The location of the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"model_name": schema.StringAttribute{
				MarkdownDescription: "The name of the model of the Node.",
				Required:            true,
			},
			"serial_number": schema.StringAttribute{
				MarkdownDescription: "The serial number of device to be associated with the Node.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"device_id": schema.StringAttribute{
				MarkdownDescription: "`device_id` defines the unique identifier of the device associated with the Node.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"position": schema.StringAttribute{
				MarkdownDescription: "The position of the Node in the Fabric.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"roles":       getRolesSchemaAttribute(),
			"metadata":    getMetadataSchemaAttribute(),
			"labels":      getLabelsSchemaAttribute(),
			"annotations": getAnnotationsSchemaAttribute(),
		},
	}
	tflog.Debug(ctx, "End schema of resource: hyperfabric_node")
}

func getRolesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of roles for a Node.`,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Set{
			setvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"LEAF", "SPINE"}...)),
		},
		ElementType: types.StringType,
	}
}

func getRolesDataSourceSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: `A set of roles for a Node.`,
		Computed:            true,
		ElementType:         types.StringType,
	}
}

func (r *NodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_node")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_node")
}

func (r *NodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_node")

	var data *NodeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_node with name '%s'", data.Name.ValueString()))

	jsonPayload := getNodeJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	result := DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/nodes", data.FabricId.ValueString()), "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client.AddChangedFabric(data.FabricId.ValueString())
	var response nodesAPIResponse
	if !decodeRestResult(&resp.Diagnostics, result, &response, "node create") {
		return
	}
	createdNode, ok := requireFirstAPIObject(&resp.Diagnostics, response.Nodes, "created node")
	if !ok {
		return
	}
	nodeID, ok := requireAPIIdentifier(&resp.Diagnostics, createdNode.NodeId, "created node")
	if !ok {
		return
	}

	data.Id = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", data.FabricId.ValueString(), nodeID))
	data.NodeId = basetypes.NewStringValue(nodeID)
	getAndSetNodeAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_node with id '%s'", data.Id.ValueString()))
}

func (r *NodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_node")
	var data *NodeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_node with id '%s'", data.Id.ValueString()))
	checkAndSetNodeIds(data)
	found := getAndSetNodeAttributes(ctx, &resp.Diagnostics, r.client, data)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_node with id '%s'", data.Id.ValueString()))
}

func (r *NodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_node")
	var data *NodeResourceModel
	var stateData *NodeResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_node with id '%s'", data.Id.ValueString()))

	jsonPayload := getNodeJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/nodes/%s", data.FabricId.ValueString(), data.NodeId.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}
	r.client.AddChangedFabric(data.FabricId.ValueString())
	getAndSetNodeAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_node with id '%s'", data.Id.ValueString()))
}

func (r *NodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_node")
	var data *NodeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_node with id '%s'", data.Id.ValueString()))
	checkAndSetNodeIds(data)
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s/nodes/%s", data.FabricId.ValueString(), data.NodeId.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client.AddChangedFabric(data.FabricId.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_node with id '%s'", data.Id.ValueString()))
}

func (r *NodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_node")
	// Processing Id in case fabricId provided is the name of the fabric and not the actual Id.
	newNode := getEmptyNodeResourceModel()
	newNode.Id = basetypes.NewStringValue(req.ID)
	checkAndSetNodeIds(newNode)
	newFabric := getEmptyFabricResourceModel()
	newFabric.Id = newNode.FabricId
	getAndSetFabricAttributes(ctx, &resp.Diagnostics, r.client, newFabric)
	newNode.FabricId = newFabric.Id
	req.ID = newNode.FabricId.ValueString() + "/nodes/" + newNode.NodeId.ValueString()
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *NodeResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_node with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_node")
}

func getAndSetNodeAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *NodeResourceModel) bool {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s/nodes/%s", data.FabricId.ValueString(), data.NodeId.ValueString()), "GET", nil)
	if diags.HasError() {
		return false
	}

	newNode := *getNewNodeResourceModelFromData(data)
	if requestData == nil || !requestData.Found {
		newNode.Id = basetypes.NewStringNull()
		*data = newNode
		return false
	}

	var response nodeAPIResponse
	if !decodeRestResult(diags, requestData, &response, "node") {
		return true
	}

	nodeID, ok := requireAPIIdentifier(diags, response.NodeId, "node")
	if !ok {
		return true
	}
	newNode.NodeId = basetypes.NewStringValue(nodeID)
	newNode.Id = basetypes.NewStringValue(fmt.Sprintf("%s/nodes/%s", newNode.FabricId.ValueString(), nodeID))
	if response.Name != nil {
		newNode.Name = basetypes.NewStringValue(*response.Name)
	}
	if response.Description != nil {
		newNode.Description = basetypes.NewStringValue(*response.Description)
	}
	if response.Enabled != nil {
		newNode.Enabled = basetypes.NewBoolValue(*response.Enabled)
	}
	if response.Location != nil {
		newNode.Location = basetypes.NewStringValue(*response.Location)
	}
	if response.ModelName != nil {
		newNode.ModelName = basetypes.NewStringValue(*response.ModelName)
	}
	if response.SerialNumber != nil {
		newNode.SerialNumber = basetypes.NewStringValue(*response.SerialNumber)
	}
	if response.DeviceId != nil {
		newNode.DeviceId = basetypes.NewStringValue(*response.DeviceId)
	}
	if response.Roles != nil {
		newNode.Roles = NewSetString(ctx, response.Roles)
	}
	if response.Metadata != nil {
		newNode.Metadata = NewMetadataObject(ctx, response.Metadata)
	}
	if response.Labels != nil {
		newNode.Labels = NewSetString(ctx, response.Labels)
	}
	if response.Annotations != nil {
		newNode.Annotations = NewNodeAnnotationsSet(ctx, response.Annotations)
		newNode.Position = NewPositionString(ctx, response.Annotations)
	}
	*data = newNode
	return true
}

func getNodeJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *NodeResourceModel, action string) *gabs.Container {
	payloadMap := map[string]interface{}{}
	payloadList := []map[string]interface{}{}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payloadMap["name"] = data.Name.ValueString()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		payloadMap["description"] = data.Description.ValueString()
	}

	// FIXME: REMOVE when PUT issue fixed
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		payloadMap["enabled"] = data.Enabled.ValueBool()
	}

	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		payloadMap["location"] = data.Location.ValueString()
	}

	if !data.ModelName.IsNull() && !data.ModelName.IsUnknown() {
		payloadMap["modelName"] = data.ModelName.ValueString()
	}

	if !data.SerialNumber.IsNull() && !data.SerialNumber.IsUnknown() {
		payloadMap["serialNumber"] = data.SerialNumber.ValueString()
	}

	if !data.DeviceId.IsNull() && !data.DeviceId.IsUnknown() {
		payloadMap["deviceId"] = data.DeviceId.ValueString()
	}

	if !data.Roles.IsNull() && !data.Roles.IsUnknown() {
		payloadMap["roles"] = getSetStringJsonPayload(ctx, data.Roles)
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		payloadMap["labels"] = getSetStringJsonPayload(ctx, data.Labels)
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		annotations := getAnnotationsJsonPayload(ctx, data.Annotations)
		if !data.Position.IsNull() && !data.Position.IsUnknown() {
			annotations = append(annotations, map[string]string{
				"name":     "position",
				"value":    data.Position.ValueString(),
				"dataType": "STRING",
			})
		}
		payloadMap["annotations"] = annotations
	}

	var payload map[string]interface{}
	if action == "create" {
		payloadList = append(payloadList, payloadMap)
		payload = map[string]interface{}{"nodes": payloadList}
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

func NewPositionString(ctx context.Context, data []interface{}) basetypes.StringValue {
	var position string
	for index, annotation := range data {
		attributes, ok := annotation.(map[string]interface{})
		if !ok {
			logUnexpectedAPIValueType(ctx, "node annotations", fmt.Sprintf("[%d]", index), "map[string]interface {}", annotation)
			continue
		}
		newAnnotation := NewAnnotationResourceModel(ctx, attributes)
		if newAnnotation.Name.ValueString() == "position" {
			position = newAnnotation.Value.ValueString()
		}
	}
	return basetypes.NewStringValue(position)
}

func checkAndSetNodeIds(data *NodeResourceModel) {
	if strings.Contains(data.Id.ValueString(), "/nodes/") {
		if data.FabricId.IsNull() || data.FabricId.IsUnknown() || data.FabricId.ValueString() == "" || data.NodeId.IsNull() || data.NodeId.IsUnknown() || data.NodeId.ValueString() == "" {
			splitId := strings.Split(data.Id.ValueString(), "/nodes/")
			data.FabricId = basetypes.NewStringValue(splitId[0])
			data.NodeId = basetypes.NewStringValue(splitId[1])
		}
	}
}
