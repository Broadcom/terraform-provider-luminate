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

func LuminateDNSServerResiliency() *schema.Resource {
	DNSServerSchema := map[string]*schema.Schema{}
	DNSServerSchema["internal_address"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "Internal address of the DNS server",
	}
	DNSServerSchema["site_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateUuid,
		Description:  "Site ID which the DNS Server will be bound for",
	}
	DNSServerSchema["name"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateString,
		Description:  "DNS server name",
	}
	DNSServerSchema["group_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: utils.ValidateUuid,
		Description:  "DNS server name",
	}
	return &schema.Resource{
		Schema:        DNSServerSchema,
		CreateContext: resourceCreateDNSResiliencyServer,
		ReadContext:   resourceReadDNSResiliencyServer,
		UpdateContext: resourceUpdateDNSResiliencyServer,
		DeleteContext: resourceDeleteDNSResiliencyServer,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateDNSResiliencyServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Creating DNS Resiliency Server")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}

	DNSGroupID := d.Get("group_id").(string)
	DNSResiliencyServerDTO := &dto.DNSServerInputDTO{
		Name:            d.Get("name").(string),
		InternalAddress: d.Get("internal_address").(string),
		SiteID:          d.Get("site_id").(string),
		GroupID:         DNSGroupID,
	}
	DNSGroup, err := client.DNSResiliencyAPI.GetDNSGroup(DNSGroupID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Reading DNS Resiliency Group with error: %s", err.Error()))
		return diag.FromErr(errors.Wrap(err, "read DNS Resiliency group failure"))
	}

	DNSResiliencyServer, err := client.DNSResiliencyAPI.CreateDNServer(DNSResiliencyServerDTO, DNSGroup.ID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Creating DNS Resiliency Server with error: %s", err.Error()))
		return diag.FromErr(errors.Wrap(err, "Create DNS Resiliency server failure"))
	}
	d.SetId(DNSResiliencyServer.ID)
	return nil
}

func resourceReadDNSResiliencyServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Reading DNS Resiliency Server")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	DNSResiliencyServerID := d.Id()
	DNSResiliencyGroupID := d.Get("group_id").(string)
	_, err := client.DNSResiliencyAPI.GetDNServer(DNSResiliencyServerID, DNSResiliencyGroupID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Reading DNS Resiliency Server with error: %s", err.Error()))
		return diag.FromErr(errors.Wrap(err, "read DNS Resiliency server failure"))
	}

	d.SetId(DNSResiliencyServerID)
	return nil
}

func resourceUpdateDNSResiliencyServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Updating DNS Resiliency Server")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}

	DNSResiliencyServerDTO := &dto.DNSServerInputDTO{
		Name:            d.Get("name").(string),
		InternalAddress: d.Get("internal_address").(string),
		SiteID:          d.Get("site_id").(string),
	}
	DNSResiliencyGroupID := d.Get("group_id").(string)
	_, err := client.DNSResiliencyAPI.UpdateDNServer(DNSResiliencyServerDTO, DNSResiliencyGroupID, d.Id())
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Updating DNS Resiliency Server with error: %s", err.Error()))
		return diag.FromErr(err)
	}
	d.SetId(d.Id())
	return nil
}

func resourceDeleteDNSResiliencyServer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[Info] Deleting DNS Resiliency Server")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}

	DNSResiliencyGroupID := d.Get("group_id").(string)
	err := client.DNSResiliencyAPI.DeleteDNSServer([]string{d.Id()}, DNSResiliencyGroupID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Deleting DNS Resiliency Server with error: %s", err.Error()))
		return diag.FromErr(errors.Wrap(err, "Delete DNS server by ids failure"))
	}
	d.SetId("")
	return nil
}
