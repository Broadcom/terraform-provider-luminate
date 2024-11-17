package provider

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateDNSGroupResiliency() *schema.Resource {
	dnsGroupSchema := map[string]*schema.Schema{}
	dnsGroupSchema["name"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "The name of DNS group",
	}
	dnsGroupSchema["send_notifications"] = &schema.Schema{
		Type:         schema.TypeBool,
		Required:     true,
		ValidateFunc: utils.ValidateBool,
		Description:  "Indicates whether notifications should be sent to admin",
	}
	dnsGroupSchema["domain_suffixes"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: utils.ValidateString,
		},
	}

	return &schema.Resource{
		Schema:        dnsGroupSchema,
		CreateContext: resourceCreateDNSResiliencyGroup,
		ReadContext:   resourceReadDNSResiliencyGroup,
		UpdateContext: resourceUpdateDNSResiliencyGroup,
		DeleteContext: resourceDeleteDNSResiliencyGroup,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateDNSResiliencyGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Creating DNS Resiliency Group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}

	DNSGroupDTO := dto.DNSGroupInputDTO{
		Name:             d.Get("name").(string),
		DomainSuffixes:   d.Get("domain_suffixes").([]interface{}),
		SendNotification: d.Get("send_notifications").(bool),
	}
	DNSGroup, err := client.DNSResiliencyAPI.CreateDNSGroup(&DNSGroupDTO)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Creating DNS Resiliency Group with error: %s", err.Error()))
		return diag.FromErr(err)
	}
	d.SetId(DNSGroup.ID)
	return resourceReadDNSResiliencyGroup(ctx, d, m)
}

func resourceReadDNSResiliencyGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Reading DNS Resiliency Group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	DNSGroupID := d.Id()
	_, err := client.DNSResiliencyAPI.GetDNSGroup(DNSGroupID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Reading DNS Resiliency Group with error: %s", err.Error()))
		return diag.FromErr(errors.Wrap(err, "read DNS Resiliency group failure"))
	}

	return nil
}

func resourceUpdateDNSResiliencyGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Updating DNS Resiliency Group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}
	DNSGroupID := d.Id()
	DNSGroup, err := client.DNSResiliencyAPI.GetDNSGroup(DNSGroupID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Reading DNS Resiliency Group with error: %s", err.Error()))
		return diag.FromErr(errors.Wrap(err, "read DNS Resiliency group failure"))
	}
	DNSGroupDTO := dto.DNSGroupInputDTO{
		Name:             d.Get("name").(string),
		DomainSuffixes:   d.Get("domain_suffixes").([]interface{}),
		SendNotification: d.Get("send_notifications").(bool),
	}
	_, err = client.DNSResiliencyAPI.UpdateDNSGroup(&DNSGroupDTO, DNSGroup.ID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Update DNS Resiliency Group with error: %s", err.Error()))
		return diag.FromErr(err)
	}
	return resourceReadDNSResiliencyGroup(ctx, d, m)
}

func resourceDeleteDNSResiliencyGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Deleting DNS Resiliency Group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	DNSGroupID := d.Id()
	err := client.DNSResiliencyAPI.DeleteDNSGroup(DNSGroupID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Deleting DNS Resiliency Group with error: %s", err.Error()))
		return diag.FromErr(errors.Wrap(err, "Delete DNS Resiliency group failure"))
	}
	d.SetId("")
	return nil
}
