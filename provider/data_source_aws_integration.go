package provider

import (
	"errors"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

func LuminateDataSourceAwsIntegration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"integration_name": {
				Type:         schema.TypeString,
				Description:  "The aws integration name as configured in Luminate portal",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"integration_id": {
				Type:        schema.TypeString,
				Description: "The aws integration id",
				Computed:    true,
			},
		},
		Read: resourceReadAwsIntegration,
	}
}

func resourceReadAwsIntegration(d *schema.ResourceData, m interface{}) error {

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	integrationName := d.Get("integration_name").(string)

	integrationId, err := client.IntegrationAPI.GetIntegrationId(integrationName)
	if err != nil {
		return err
	}

	d.SetId(integrationId)

	return nil
}
