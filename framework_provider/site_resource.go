package framework_provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LuminateSiteResource struct {
	BaseLuminateResource
}

type LuminateSiteResourceModel struct {
	ID                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Region                         types.String `tfsdk:"region"`
	MuteHealthNotifications        types.Bool   `tfsdk:"mute_health_notification"`
	KubernetesPersistentVolumeName types.String `tfsdk:"kubernetes_persistent_volume_name"`
}

func NewLuminateSiteResource() resource.Resource {
	return &LuminateSiteResource{}
}

func (r *LuminateSiteResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_site2"
}

func (r *LuminateSiteResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Site name",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Site name",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Site connectivity region",
			},
			"mute_health_notification": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Mute notifications if site is down",
				Default:     booldefault.StaticBool(false),
			},
			"kubernetes_persistent_volume_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Kubernetes persistent volume name",
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func siteModelToDto(data *LuminateSiteResourceModel) *dto.Site {
	site := dto.Site{
		Name:       data.Name.ValueString(),
		Region:     data.Region.ValueString(),
		MuteHealth: data.MuteHealthNotifications.ValueBool(),
		K8SVolume:  data.KubernetesPersistentVolumeName.ValueString(),
	}

	return &site
}

func siteDtoToModel(siteDto *dto.Site, data *LuminateSiteResourceModel) {
	data.ID = types.StringValue(siteDto.ID)
	data.Name = types.StringValue(siteDto.Name)
	data.Region = types.StringValue(siteDto.Region)
	data.MuteHealthNotifications = types.BoolValue(siteDto.MuteHealth)
	data.KubernetesPersistentVolumeName = types.StringValue(siteDto.K8SVolume)
}

func (r *LuminateSiteResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data LuminateSiteResourceModel

	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	site := siteModelToDto(&data)
	newSite, err := r.client.Sites.CreateSite(site)
	if err != nil {
		response.Diagnostics.AddError("Failed to create site", err.Error())
		return
	}

	siteDtoToModel(newSite, &data)
	diag := response.State.Set(ctx, &data)
	response.Diagnostics.Append(diag...)
}

func (r *LuminateSiteResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data LuminateSiteResourceModel

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	actualSite, err := r.client.Sites.GetSiteByID(data.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to get site", err.Error())
		return
	}

	if actualSite == nil {
		response.State.RemoveResource(ctx)
		return
	}

	siteDtoToModel(actualSite, &data)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *LuminateSiteResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var currentData LuminateSiteResourceModel
	var data LuminateSiteResourceModel

	response.Diagnostics.Append(request.State.Get(ctx, &currentData)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	site := siteModelToDto(&data)
	updatedSite, err := r.client.Sites.UpdateSite(site, currentData.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to create site", err.Error())
		return
	}

	siteDtoToModel(updatedSite, &data)
	data.ID = currentData.ID
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *LuminateSiteResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data LuminateSiteResourceModel

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.Sites.DeleteSite(data.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to delete site", err.Error())
		return
	}
}
