package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	swagger "github.gwd.broadcom.net/SED/ztna-api-documentation/go/sdk"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func LuminateAccessPolicyBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:         schema.TypeBool,
			Description:  "Indicates whether this policy is enabled.",
			Optional:     true,
			Default:      true,
			ValidateFunc: utils.ValidateBool,
		},
		"name": {
			Type:         schema.TypeString,
			Description:  "A descriptive name of the policy.",
			Required:     true,
			ValidateFunc: utils.ValidatePolicyName,
		},
		"identity_provider_id": {
			Type:         schema.TypeString,
			Description:  "The identity provider id",
			Optional:     true,
			ValidateFunc: utils.ValidateString,
			ForceNew:     true,
		},
		"user_ids": {
			Type:        schema.TypeList,
			Description: "The user entities to which this policy applies.",
			Optional:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: utils.ValidateString,
			},
		},
		"group_ids": {
			Type:        schema.TypeList,
			Description: "The group entities to which this policy applies.",
			Optional:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: utils.ValidateString,
			},
		},
		"applications": {
			Type:        schema.TypeSet,
			Description: "The applications to which this policy applies.",
			Required:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: utils.ValidateUuid,
			},
		},
		"collection_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Collection ID to which the policy will be assigned",
		},
		"validators": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"web_verification": {
						Type:         schema.TypeBool,
						Optional:     true,
						Default:      false,
						Description:  "Indicate whatever to perform web verification validation. not compatible for HTTP applications",
						ValidateFunc: utils.ValidateBool,
					},
				},
			},
		},
		"conditions": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"location": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "location based condition, specify the list of accepted locations.",
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.NoZeroValues,
						},
					},
					"source_ip": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "source ip based condition, specify the allowed CIDR for this policy.",
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.NoZeroValues,
						},
					},
					"managed_device": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "list of managed devices that have restriction access",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"opswat": {
									Type:         schema.TypeBool,
									Optional:     true,
									Default:      false,
									Description:  "Indicate whatever to restrict access to Opswat MetaAccess",
									ValidateFunc: utils.ValidateBool,
								},
								"symantec_cloudsoc": {
									Type:         schema.TypeBool,
									Optional:     true,
									Default:      false,
									Description:  "Indicate whatever to restrict access to symantec cloudsoc",
									ValidateFunc: utils.ValidateBool,
								},
								"symantec_web_security_service": {
									Type:         schema.TypeBool,
									Optional:     true,
									Default:      false,
									Description:  "Indicate whatever to restrict access to symantec web security service",
									ValidateFunc: utils.ValidateBool,
								},
							},
						},
					},
					"unmanaged_device": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "list of unmanaged devices that have restriction access",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"opswat": {
									Type:         schema.TypeBool,
									Optional:     true,
									Default:      false,
									Description:  "Indicate whatever to restrict access to Opswat MetaAccess",
									ValidateFunc: utils.ValidateBool,
								},
								"symantec_cloudsoc": {
									Type:         schema.TypeBool,
									Optional:     true,
									Default:      false,
									Description:  "Indicate whatever to restrict access to symantec cloudsoc",
									ValidateFunc: utils.ValidateBool,
								},
								"symantec_web_security_service": {
									Type:         schema.TypeBool,
									Optional:     true,
									Default:      false,
									Description:  "Indicate whatever to restrict access to symantec web security service",
									ValidateFunc: utils.ValidateBool,
								},
							},
						},
					},
				},
			},
		},
	}
}

func setAccessPolicyBaseFields(d *schema.ResourceData, accessPolicy *dto.AccessPolicy) {
	d.SetId(accessPolicy.Id)
	d.Set("enabled", accessPolicy.Enabled)
	d.Set("name", accessPolicy.Name)
	d.Set("applications", accessPolicy.Applications)
	d.Set("collection_id", accessPolicy.CollectionID)

	if accessPolicy.Validators != nil {
		d.Set("validators", flattenValidators(accessPolicy.Validators))
	}

	if accessPolicy.Conditions != nil {
		d.Set("conditions", flattenConditions(accessPolicy.Conditions))
	}
	if len(accessPolicy.DirectoryEntities) > 0 {
		groupIDs := make([]string, 0)
		userIDs := make([]string, 0)
		for _, entity := range accessPolicy.DirectoryEntities {
			if *dto.ToModelType(entity.EntityType) == swagger.GROUP_EntityType {
				groupIDs = append(groupIDs, entity.IdentifierInProvider)
			}
			if *dto.ToModelType(entity.EntityType) == swagger.USER_EntityType {
				userIDs = append(userIDs, entity.IdentifierInProvider)
			}
		}
		d.Set("group_ids", groupIDs)
		d.Set("user_ids", userIDs)
	}
}

func flattenValidators(validators *dto.Validators) []interface{} {
	k := make(map[string]interface{})
	if validators == nil {
		return []interface{}{}
	}
	if validators.WebVerification {
		k["web_verification"] = true
	}
	if validators.MFA {
		k["mfa"] = true
	}
	return []interface{}{k}
}

func flattenManagedDevice(manageDevice dto.Device) []interface{} {
	var out = make([]interface{}, 0, 0)
	k := make(map[string]interface{})

	if manageDevice.OpswatMetaAccess {
		k["opswat"] = manageDevice.OpswatMetaAccess
	}

	if manageDevice.SymantecCloudSoc {
		k["symantec_cloudsoc"] = manageDevice.SymantecCloudSoc
	}

	if manageDevice.SymantecWebSecurityService {
		k["symantec_web_security_service"] = manageDevice.SymantecWebSecurityService
	}
	if len(k) == 0 {
		return nil
	}
	out = append(out, k)
	return out
}

func flattenConditions(conditions *dto.Conditions) []interface{} {
	if conditions == nil {
		return []interface{}{}
	}

	k := map[string]interface{}{
		"source_ip": conditions.SourceIp,
		"location":  conditions.Location,
	}

	if hasDeviceCondition(conditions.ManagedDevice) {
		k["managed_device"] = flattenManagedDevice(conditions.ManagedDevice)
	}

	if hasDeviceCondition(conditions.ManagedDevice) {
		k["unmanaged_device"] = flattenManagedDevice(conditions.UnmanagedDevice)
	}

	return []interface{}{k}
}

func hasDeviceCondition(managedDevice dto.Device) bool {
	if managedDevice.OpswatMetaAccess {
		return true
	}

	if managedDevice.SymantecCloudSoc {
		return true
	}

	if managedDevice.SymantecWebSecurityService {
		return true
	}
	return false
}

func extractAccessPolicyBaseFields(d *schema.ResourceData) *dto.AccessPolicy {
	var applicationIds []string
	var directoryEntity []dto.DirectoryEntity
	var validators *dto.Validators
	var conditions *dto.Conditions

	enabled := d.Get("enabled").(bool)
	name := d.Get("name").(string)
	identityProviderId := d.Get("identity_provider_id").(string)

	userIdsInterface := d.Get("user_ids").([]interface{})

	collectionID := d.Get("collection_id").(string)

	for _, userId := range userIdsInterface {
		directoryEntity = append(directoryEntity, dto.DirectoryEntity{
			IdentityProviderId:   identityProviderId,
			IdentifierInProvider: userId.(string),
			EntityType:           "User",
		})
	}

	groupIdsInterface := d.Get("group_ids").([]interface{})

	for _, groupId := range groupIdsInterface {
		directoryEntity = append(directoryEntity, dto.DirectoryEntity{
			IdentityProviderId:   identityProviderId,
			IdentifierInProvider: groupId.(string),
			EntityType:           "Group",
		})
	}

	applicationIdsInterface := d.Get("applications").(*schema.Set)
	applicationIdsList := applicationIdsInterface.List()
	for _, applicationId := range applicationIdsList {
		applicationIds = append(applicationIds, applicationId.(string))
	}
	validators = extractValidators(d)
	conditions = extractConditions(d)

	return &dto.AccessPolicy{
		Policy: dto.Policy{
			Enabled:           enabled,
			Name:              name,
			DirectoryEntities: directoryEntity,
			Applications:      applicationIds,
			Conditions:        conditions,
			CollectionID:      collectionID,
		},
		Validators: validators,
	}
}

func extractValidators(d *schema.ResourceData) *dto.Validators {
	var validators *dto.Validators

	if v, ok := d.GetOk("validators"); ok {
		var webVerification bool
		var mfaVerification bool

		for _, element := range v.([]interface{}) {
			elem := element.(map[string]interface{})
			if value, ok := elem["web_verification"].(bool); ok && value {
				webVerification = true
			}
			if value, ok := elem["mfa"].(bool); ok && value {
				mfaVerification = true
			}
		}

		if webVerification {
			validators = &dto.Validators{
				WebVerification: webVerification,
			}
		}

		if mfaVerification {
			validators = &dto.Validators{
				MFA: true,
			}
		}
	}

	return validators
}

func extractConditions(d *schema.ResourceData) *dto.Conditions {
	var conditions *dto.Conditions

	if v, ok := d.GetOk("conditions"); ok {
		for _, element := range v.([]interface{}) {
			elem := element.(map[string]interface{})

			var sourceIpList []string
			var locations []string
			var managedDevice dto.Device
			var unmanagedDevice dto.Device

			if sourceIpInterface, ok := elem["source_ip"].([]interface{}); ok {
				for _, sourceIp := range sourceIpInterface {
					sourceIpList = append(sourceIpList, sourceIp.(string))
				}
			}

			if locationsInterface, ok := elem["location"].([]interface{}); ok {
				for _, location := range locationsInterface {
					locations = append(locations, location.(string))
				}
			}

			if managedDeviceInterface, ok := elem["managed_device"].([]interface{}); ok {
				deviceList(managedDeviceInterface, &managedDevice)
			}

			if unManagedDeviceInterface, ok := elem["unmanaged_device"].([]interface{}); ok {
				deviceList(unManagedDeviceInterface, &unmanagedDevice)
			}

			conditions = &dto.Conditions{
				SourceIp:        sourceIpList,
				Location:        locations,
				ManagedDevice:   managedDevice,
				UnmanagedDevice: unmanagedDevice,
			}
		}
	}

	return conditions
}

func deviceList(deviceInterface []interface{}, device *dto.Device) {
	for _, deviceElements := range deviceInterface {
		elem := deviceElements.(map[string]interface{})

		if elem["opswat"].(bool) {
			device.OpswatMetaAccess = elem["opswat"].(bool)
		}

		if elem["symantec_cloudsoc"].(bool) {
			device.SymantecCloudSoc = elem["symantec_cloudsoc"].(bool)
		}

		if elem["symantec_web_security_service"].(bool) {
			device.SymantecWebSecurityService = elem["symantec_web_security_service"].(bool)
		}
	}
}

func resourceDeleteAccessPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	err := client.AccessPolicies.DeleteAccessPolicy(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
