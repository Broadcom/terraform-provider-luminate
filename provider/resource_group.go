package provider

import (
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group Name",
				ForceNew:    true,
			},
			"identity_provider_id": {
				Type:        schema.TypeString,
				Description: "The identity provider id",
				Required:    true,
				ForceNew:    true,
			},
		},

		CreateContext: resourceCreateGroup,
		ReadContext:   resourceReadGroup,
		DeleteContext: resourceDeleteGroup,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] creating group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}

	groupName := d.Get("name").(string)
	idpID := d.Get("identity_provider_id").(string)
	group, err := client.Groups.CreateGroup(idpID, groupName)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Creating Group %s in idp %s with error: %s", groupName, idpID, err.Error()))
		return diag.FromErr(err)
	}
	d.SetId(group.ID)
	d.Set("identity_provider_id", group.IdentityProviderId)
	d.Set("name", group.Name)
	return resourceReadGroup(ctx, d, m)
}

func resourceReadGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading group")
	idpID := d.Get("identity_provider_id").(string)
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	groupID := d.Id()
	group, err := client.Groups.GetGroup(groupID, idpID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed Reading Group %s in idp %s with error: %s", groupID, idpID, err.Error()))
		return diag.FromErr(err)
	}
	d.Set("identity_provider_id", group.IdentityProviderId)
	d.Set("name", group.Name)
	return nil
}

func resourceDeleteGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}

	groupID := d.Id()
	idpID := d.Get("identity_provider_id").(string)
	err := client.Groups.DeleteGroup(idpID, groupID)
	if err != nil {
		log.Println(fmt.Sprintf("[Error] failed deleting group %s in idp %s with error: %s", groupID, idpID, err.Error()))
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
