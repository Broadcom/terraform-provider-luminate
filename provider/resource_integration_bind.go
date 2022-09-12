package provider

import (
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

func LuminateawsIntegrationBind() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"integration_name": {
				Type:         schema.TypeString,
				Description:  "The aws integration name as configured in Luminate portal",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"integration_id": {
				Type:         schema.TypeString,
				Description:  "The aws integration id",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"aws_role_arn": {
				Type:         schema.TypeString,
				Description:  "The aws role arn",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"luminate_aws_account_id": {
				Type:         schema.TypeString,
				Description:  "The Luminate aws account id",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"aws_external_id": {
				Type:         schema.TypeString,
				Description:  "The aws external id",
				Required:     true,
				ValidateFunc: utils.ValidateString,
			},
			"regions": {
				Type:        schema.TypeList,
				Description: "regions",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Create: resourceCreateAwsIntegrationBind,
		Read:   resourceReadAwsIntegrationBind,
		Update: resourceUpdateIntegrationBind,
		Delete: resourceDeleteAwsIntegrationBind,
	}
}

func resourceCreateAwsIntegrationBind(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE INTEGRATION BIND CREATE")

	return resourceUpdateIntegrationBind(d, m)
}
func resourceReadAwsIntegrationBind(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE INTEGRATION BIND READ")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	integrationID := d.Get("integration_id").(string)

	awsIntegration, err := client.IntegrationAPI.ReadAWSIntegrationBind(integrationID)
	if err != nil {
		return err
	}

	setAwsIntegrationFieldsBind(d, awsIntegration)

	return nil
}

func resourceUpdateIntegrationBind(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE INTEGRATION BIND UPDATE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	awsBody := createAWSrequestBody(d)

	_, err := client.IntegrationAPI.UpdateAWSIntegration(awsBody)
	if err != nil {
		return err
	}
	d.SetId(awsBody.ID)

	return nil
}

func resourceDeleteAwsIntegrationBind(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE INTEGRATION BIND DELETE")

	return nil
}

func setAwsIntegrationFieldsBind(d *schema.ResourceData, integration *dto.AwsIntegrationBind) {
	d.Set("integration_name", integration.Name)
	d.Set("integration_id", integration.Id)
	d.Set("luminate_aws_account_id", integration.LuminateAwsAccountId)
	d.Set("aws_external_id", integration.AwsExternalId)
	d.Set("aws_role_arn", integration.AwsRoleArn)
}

func createAWSrequestBody(d *schema.ResourceData) *service.AWSRequestBody {
	integrationName := d.Get("integration_name").(string)
	integrationID := d.Get("integration_id").(string)
	awsArn := d.Get("aws_role_arn").(string)
	luminateAwsID := d.Get("luminate_aws_account_id").(string)
	awsExternalID := d.Get("aws_external_id").(string)
	regionsInterface := d.Get("regions").([]interface{})

	regions := make([]string, 0, len(regionsInterface))
	for i := range regionsInterface {
		region, _ := regionsInterface[i].(string)
		regions = append(regions, region)
	}

	req := service.AWSRequestBody{
		Provider:             "amazon",
		HostnameTagName:      "Name",
		Name:                 integrationName,
		AwsExternalID:        awsExternalID,
		ID:                   integrationID,
		LuminateAwsAccountID: luminateAwsID,
		Regions:              regions,
		AwsRoleArn:           awsArn,
	}
	return &req
}
