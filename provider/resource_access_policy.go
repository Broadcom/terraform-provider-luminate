package provider

import (
	"errors"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			Type:        schema.TypeList,
			Description: "The applications to which this policy applies.",
			Required:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: utils.ValidateUuid,
			},
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

	if accessPolicy.Validators != nil {
		d.Set("validators", flattenValidators(accessPolicy.Validators))
	}

	if accessPolicy.Conditions != nil {
		d.Set("conditions", flattenConditions(accessPolicy.Conditions))
	}
}

func flattenValidators(validators *dto.Validators) []interface{} {
	if validators == nil {
		return []interface{}{}
	}
	k := map[string]interface{}{
		"compliance_check": validators.ComplianceCheck,
		"web_verification": validators.WebVerification,
	}
	return []interface{}{k}
}

func flattenManagedDevice(manageDevice dto.ManagedDevice) []interface{} {
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

	out = append(out, k)
	return out
}

func flattenConditions(conditions *dto.Conditions) []interface{} {
	if conditions == nil {
		return []interface{}{}
	}

	k := map[string]interface{}{
		"source_ip":        conditions.SourceIp,
		"location":         conditions.Location,
		"managed_device":   flattenManagedDevice(conditions.ManagedDevice),
		"unmanaged_device": conditions.UnmanagedDevice,
	}

	return []interface{}{k}
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

	applicationIdsInterface := d.Get("applications").([]interface{})
	for _, applicationId := range applicationIdsInterface {
		applicationIds = append(applicationIds, applicationId.(string))
	}

	validators = extractValidators(d)
	conditions = extractConditions(d)

	return &dto.AccessPolicy{
		Enabled:           enabled,
		Name:              name,
		DirectoryEntities: directoryEntity,
		Applications:      applicationIds,
		Validators:        validators,
		Conditions:        conditions,
	}
}

func extractValidators(d *schema.ResourceData) *dto.Validators {
	var validators *dto.Validators

	if v, ok := d.GetOk("validators"); ok {
		var complianceCheck bool
		var webVerification bool

		for _, element := range v.([]interface{}) {
			elem := element.(map[string]interface{})
			if value, ok := elem["compliance_check"].(bool); ok && value {
				complianceCheck = true
			}

			if value, ok := elem["web_verification"].(bool); ok && value {
				webVerification = true
			}
		}

		if complianceCheck || webVerification {
			validators = &dto.Validators{
				ComplianceCheck: complianceCheck,
				WebVerification: webVerification,
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
			var managedDevice dto.ManagedDevice
			var unmanagedDevice bool

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
				for _, managedDeviceElements := range managedDeviceInterface {
					elem := managedDeviceElements.(map[string]interface{})

					if elem["opswat"].(bool) {
						managedDevice.OpswatMetaAccess = elem["opswat"].(bool)
					}

					if elem["symantec_cloudsoc"].(bool) {
						managedDevice.SymantecCloudSoc = elem["symantec_cloudsoc"].(bool)
					}

					if elem["symantec_web_security_service"].(bool) {
						managedDevice.SymantecWebSecurityService = elem["symantec_web_security_service"].(bool)
					}
				}
			}

			unmanagedDevice, ok = elem["unmanaged_device"].(bool)
			if !ok {
				unmanagedDevice = false
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

func resourceDeleteAccessPolicy(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	err := client.AccessPolicies.DeleteAccessPolicy(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
