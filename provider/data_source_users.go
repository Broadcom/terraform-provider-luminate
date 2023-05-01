package provider

import (
	"errors"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func LuminateDataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identity_provider_id": {
				Type:         schema.TypeString,
				Description:  "The identity provider id",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"users": {
				Type:        schema.TypeList,
				Description: "list of users to include as part of this policy",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: utils.ValidateEmail,
				},
			},
			"user_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Computed: true,
			},
		},
		Read: resourceReadUsers,
	}
}

func resourceReadUsers(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	var userIds []string

	identityProviderId := d.Get("identity_provider_id").(string)
	userNames := d.Get("users").([]interface{})

	for _, userEmail := range userNames {
		userID, err := client.Users.GetUserId(identityProviderId, userEmail.(string))
		if err != nil {
			return err
		}

		userIds = append(userIds, userID)
	}

	d.SetId(identityProviderId)
	d.Set("user_ids", userIds)

	return nil
}
