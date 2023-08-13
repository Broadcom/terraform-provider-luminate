package provider

import (
	"context"
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LuminateAWSIntegration() *schema.Resource {
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
			"luminate_aws_account_id": {
				Type:        schema.TypeString,
				Description: "The Luminate aws account id",
				Computed:    true,
			},
			"aws_external_id": {
				Type:        schema.TypeString,
				Description: "The aws external id",
				Computed:    true,
			},
		},
		CreateContext: resourceCreateAwsIntegration,
		ReadContext:   resourceAwsReadIntegration,
		UpdateContext: resourceUpdateAwsIntegration,
		DeleteContext: resourceDeleteAwsIntegration,
	}
}

func resourceCreateAwsIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE INTEGRATION - CREATE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	integrationName := d.Get("integration_name").(string)

	newIntegration, err := client.IntegrationAPI.CreateAWSIntegration(integrationName)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(newIntegration.Id)
	setAwsIntegrationFields(d, newIntegration)

	return resourceAwsReadIntegration(ctx, d, m)
}

func resourceAwsReadIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE INTEGRATION - READ")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	integrationID := d.Get("integration_id").(string)

	awsIntegration, err := client.IntegrationAPI.ReadAWSIntegration(integrationID)
	if err != nil {
		return diag.FromErr(err)
	}

	setAwsIntegrationFields(d, awsIntegration)

	return diagnostics
}

func resourceUpdateAwsIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE INTEGRATION - UPDATE")

	return resourceAwsReadIntegration(ctx, d, m)
}

func resourceDeleteAwsIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE INTEGRATION - DELETE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	integrationID := d.Get("integration_id").(string)

	err := client.IntegrationAPI.DeleteAWSIntegration(integrationID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setAwsIntegrationFields(d *schema.ResourceData, integration *dto.AwsIntegration) {
	d.Set("integration_id", integration.Id)
	d.Set("luminate_aws_account_id", integration.LuminateAwsAccountId)
	d.Set("aws_external_id", integration.AwsExternalId)
}
