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

	"github.com/CiscoDevNet/terraform-provider-hyperfabric/internal/client"
	"github.com/Jeffail/gabs/v2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FabricResource{}
var _ resource.ResourceWithImportState = &FabricResource{}

func NewFabricResource() resource.Resource {
	return &FabricResource{}
}

// FabricResource defines the resource implementation.
type FabricResource struct {
	client *client.Client
}

// FabricResourceModel describes the resource data model.
type FabricResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	// Enabled     types.Bool   `tfsdk:"enabled"`
	Topology    types.String `tfsdk:"topology"`
	Location    types.String `tfsdk:"location"`
	Address     types.String `tfsdk:"address"`
	City        types.String `tfsdk:"city"`
	Country     types.String `tfsdk:"country"`
	Metadata    types.Object `tfsdk:"metadata"`
	Labels      types.Set    `tfsdk:"labels"`
	Annotations types.Set    `tfsdk:"annotations"`
}

func getEmptyFabricResourceModel() *FabricResourceModel {
	return &FabricResourceModel{
		Id:          basetypes.NewStringNull(),
		Name:        basetypes.NewStringNull(),
		Description: basetypes.NewStringNull(),
		// Enabled:     basetypes.NewBoolValue(true),
		Topology:    basetypes.NewStringNull(),
		Location:    basetypes.NewStringNull(),
		Address:     basetypes.NewStringNull(),
		City:        basetypes.NewStringNull(),
		Country:     basetypes.NewStringNull(),
		Metadata:    basetypes.NewObjectNull(MetadataResourceModelAttributeType()),
		Labels:      basetypes.NewSetNull(SetStringResourceModelAttributeType()),
		Annotations: basetypes.NewSetNull(AnnotationResourceModelAttributeType()),
	}
}

func getNewFabricResourceModelFromData(data *FabricResourceModel) *FabricResourceModel {
	newFabric := getEmptyFabricResourceModel()

	if !data.Id.IsNull() && !data.Id.IsUnknown() {
		newFabric.Id = data.Id
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		newFabric.Name = data.Name
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		newFabric.Description = data.Description
	}

	// if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
	//  newFabric.Enabled = data.Enabled
	// }

	if !data.Topology.IsNull() && !data.Topology.IsUnknown() {
		newFabric.Topology = data.Topology
	}

	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		newFabric.Location = data.Location
	}

	if !data.Address.IsNull() && !data.Address.IsUnknown() {
		newFabric.Address = data.Address
	}

	if !data.City.IsNull() && !data.City.IsUnknown() {
		newFabric.City = data.City
	}

	if !data.Country.IsNull() && !data.Country.IsUnknown() {
		newFabric.Country = data.Country
	}

	if !data.Labels.IsNull() && !data.Labels.IsUnknown() {
		newFabric.Labels = data.Labels
	}

	if !data.Annotations.IsNull() && !data.Annotations.IsUnknown() {
		newFabric.Annotations = data.Annotations
	}
	return newFabric
}

type FabricIdentifier struct {
	Id types.String
}

func (r *FabricResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "Start metadata of resource: hyperfabric_fabric")
	resp.TypeName = req.ProviderTypeName + "_fabric"
	tflog.Debug(ctx, "End metadata of resource: hyperfabric_fabric")
}

func (r *FabricResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "Start schema of resource: hyperfabric_fabric")
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`id` defines the unique identifier of a Fabric.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Fabric.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description is a user defined field to store notes about the Fabric.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"topology": schema.StringAttribute{
				MarkdownDescription: "The topology type of the Fabric.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MESH", "SPINE_LEAF"}...),
				},
			},
			// "enabled": schema.BoolAttribute{
			// 	MarkdownDescription: "The enabled state of the Fabric.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	PlanModifiers: []planmodifier.Bool{
			// 		boolplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			"location": schema.StringAttribute{
				MarkdownDescription: "The location of the Fabric.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address where the Fabric is located.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "The city where the Fabric is located.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					SetToStringNullWhenStateIsNullPlanIsUnknownDuringUpdate(),
				},
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "The country in which the Fabric is located.",
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
	tflog.Debug(ctx, "End schema of resource: hyperfabric_fabric")
}

func (r *FabricResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Start configure of resource: hyperfabric_fabric")
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
	tflog.Debug(ctx, "End configure of resource: hyperfabric_fabric")
}

func (r *FabricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start create of resource: hyperfabric_fabric")

	var data *FabricResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Create of resource hyperfabric_fabric with name '%s'", data.Name.ValueString()))

	jsonPayload := getFabricJsonPayload(ctx, &resp.Diagnostics, data, "create")
	if resp.Diagnostics.HasError() {
		return
	}

	result := DoRestRequest(ctx, &resp.Diagnostics, r.client, "/api/v1/fabrics", "POST", jsonPayload)
	if resp.Diagnostics.HasError() {
		return
	}

	var response fabricsAPIResponse
	if !decodeRestResult(&resp.Diagnostics, result, &response, "fabric create") {
		return
	}
	createdFabric, ok := requireFirstAPIObject(&resp.Diagnostics, response.Fabrics, "created fabric")
	if !ok {
		return
	}
	fabricID, ok := requireAPIIdentifier(&resp.Diagnostics, createdFabric.FabricId, "created fabric")
	if !ok {
		return
	}

	data.Id = basetypes.NewStringValue(fabricID)
	getAndSetFabricAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End create of resource hyperfabric_fabric with name '%s'", data.Name.ValueString()))
}

func (r *FabricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start read of resource: hyperfabric_fabric")
	var data *FabricResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read of resource hyperfabric_fabric with id '%s'", data.Id.ValueString()))

	found := getAndSetFabricAttributes(ctx, &resp.Diagnostics, r.client, data)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Debug(ctx, fmt.Sprintf("End read of resource hyperfabric_fabric with id '%s'", data.Id.ValueString()))
}

func (r *FabricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start update of resource: hyperfabric_fabric")
	var data *FabricResourceModel
	var stateData *FabricResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update of resource hyperfabric_fabric with id '%s'", data.Id.ValueString()))

	jsonPayload := getFabricJsonPayload(ctx, &resp.Diagnostics, data, "update")

	if resp.Diagnostics.HasError() {
		return
	}

	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s", data.Id.ValueString()), "PUT", jsonPayload)

	if resp.Diagnostics.HasError() {
		return
	}

	getAndSetFabricAttributes(ctx, &resp.Diagnostics, r.client, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("End update of resource hyperfabric_fabric with id '%s'", data.Id.ValueString()))
}

func (r *FabricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start delete of resource: hyperfabric_fabric")
	var data *FabricResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete of resource hyperfabric_fabric with id '%s'", data.Id.ValueString()))
	DoRestRequest(ctx, &resp.Diagnostics, r.client, fmt.Sprintf("/api/v1/fabrics/%s", data.Id.ValueString()), "DELETE", nil)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("End delete of resource hyperfabric_fabric with id '%s'", data.Id.ValueString()))
}

func (r *FabricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start import state of resource: hyperfabric_fabric")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	var stateData *FabricResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	tflog.Debug(ctx, fmt.Sprintf("Import state of resource hyperfabric_fabric with id '%s'", stateData.Id.ValueString()))
	tflog.Debug(ctx, "End import of state resource: hyperfabric_fabric")
}

func getAndSetFabricAttributes(ctx context.Context, diags *diag.Diagnostics, client *client.Client, data *FabricResourceModel) bool {
	requestData := DoRestRequest(ctx, diags, client, fmt.Sprintf("/api/v1/fabrics/%s", data.Id.ValueString()), "GET", nil)
	if diags.HasError() {
		return false
	}

	newFabric := *getNewFabricResourceModelFromData(data)
	if requestData == nil || !requestData.Found {
		newFabric.Id = basetypes.NewStringNull()
		*data = newFabric
		return false
	}

	var response fabricAPIResponse
	if !decodeRestResult(diags, requestData, &response, "fabric") {
		return true
	}

	fabricID, ok := requireAPIIdentifier(diags, response.FabricId, "fabric")
	if !ok {
		return true
	}
	newFabric.Id = basetypes.NewStringValue(fabricID)
	if response.Name != nil {
		newFabric.Name = basetypes.NewStringValue(*response.Name)
	}
	if response.Description != nil {
		newFabric.Description = basetypes.NewStringValue(*response.Description)
	}
	if response.Topology != nil {
		newFabric.Topology = basetypes.NewStringValue(*response.Topology)
	}
	if response.Location != nil {
		newFabric.Location = basetypes.NewStringValue(*response.Location)
	}
	if response.Address != nil {
		newFabric.Address = basetypes.NewStringValue(*response.Address)
	}
	if response.City != nil {
		newFabric.City = basetypes.NewStringValue(*response.City)
	}
	if response.Country != nil {
		newFabric.Country = basetypes.NewStringValue(*response.Country)
	}
	if response.Metadata != nil {
		newFabric.Metadata = NewMetadataObject(ctx, response.Metadata)
	}
	if response.Labels != nil {
		newFabric.Labels = NewSetString(ctx, response.Labels)
	}
	if response.Annotations != nil {
		newFabric.Annotations = NewAnnotationsSet(ctx, response.Annotations)
	}
	*data = newFabric
	return true
}

func getFabricJsonPayload(ctx context.Context, diags *diag.Diagnostics, data *FabricResourceModel, action string) *gabs.Container {
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

	if !data.Topology.IsNull() && !data.Topology.IsUnknown() {
		payloadMap["topology"] = data.Topology.ValueString()
	}

	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		payloadMap["location"] = data.Location.ValueString()
	}

	if !data.Address.IsNull() && !data.Address.IsUnknown() {
		payloadMap["address"] = data.Address.ValueString()
	}

	if !data.City.IsNull() && !data.City.IsUnknown() {
		payloadMap["city"] = data.City.ValueString()
	}

	if !data.Country.IsNull() && !data.Country.IsUnknown() {
		payloadMap["country"] = data.Country.ValueString()
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
		payload = map[string]interface{}{"fabrics": payloadList}
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
