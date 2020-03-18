package provider

import (
	"errors"
	"fmt"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

func LuminateConnector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "A descriptive name of the Connector",
				Required:     true,
				ValidateFunc: utils.ValidateString,
				ForceNew:     true,
			},
			"site_id": {
				Type:         schema.TypeString,
				Description:  "Site to bind the connector to",
				Required:     true,
				ValidateFunc: utils.ValidateUuid,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "The deployment type of the host running the Luminate connector",
				Required:     true,
				ValidateFunc: validateConnectorType,
				ForceNew:     true,
			},
			"command": {
				Type:        schema.TypeString,
				Description: "Command for deploying Luminate connector",
				Computed:    true,
			},
			"otp": {
				Type:        schema.TypeString,
				Description: "One time password for running Luminate connector",
				Computed:    true,
			},
		},
		Create: resourceCreateConnector,
		Read:   resourceReadConnector,
		Delete: resourceDeleteConnector,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateConnector(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	conOpts := extractConnectorFields(d)

	newCon, err := client.Connectors.CreateConnector(conOpts, conOpts.SiteID)
	if err != nil {
		return err
	}

	setConnectorFields(d, newCon)

	command, err := client.Connectors.GetConnectorCommand(newCon.ID)
	if err != nil {
		return err
	}

	d.Set("command", command)

	return resourceReadConnector(d, m)
}

func resourceReadConnector(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	connector, err := client.Connectors.GetConnectorByID(d.Id())
	if err != nil {
		return err
	}

	if connector == nil {
		d.SetId("")
		return nil
	}

	connector.SiteID = d.Get("site_id").(string)

	setConnectorFields(d, connector)

	return nil
}

func resourceDeleteConnector(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	err := client.Connectors.DeleteConnector(d.Id())
	if err != nil {
		return errors.New("unable to delete connector")
	}
	d.SetId("")

	return resourceReadConnector(d, m)
}

func validateConnectorType(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	cType, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type to be string"))
		return warns, errs
	}

	validTypes := []string{
		"linux",
		"kubernetes",
		"windows",
		"docker-compose",
	}

	if !utils.StringInSlice(validTypes, cType) {
		errs = append(errs, fmt.Errorf("connector type must be one of %v", validTypes))
	}
	return warns, errs
}

func setConnectorFields(d *schema.ResourceData, connector *dto.Connector) {
	d.SetId(connector.ID)
	d.Set("name", connector.Name)
	d.Set("type", connector.Type)
	d.Set("site_id", connector.SiteID)
	d.Set("otp", connector.OTP)
}

func extractConnectorFields(d *schema.ResourceData) *dto.Connector {
	return &dto.Connector{
		Name:   d.Get("name").(string),
		Type:   d.Get("type").(string),
		SiteID: d.Get("site_id").(string),
	}
}
