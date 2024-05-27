package provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"log"
)

const IDPID = "9748fcbe-7eac-4dc9-9809-127ca4f889ba"

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
				Description: "list of users to include as part of this policy",
				Computed:    true,
			},
		},

		CreateContext: resourceCreateGroup,
		ReadContext:   resourceReadGroup,
		//UpdateContext: resourceUpdateGroup,
		DeleteContext: resourceDeleteGroup,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}
	groupName := d.Get("name").(string)
	group, err := client.Groups.CreateGroup(IDPID, groupName)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to create group"))
	}
	d.SetId(group.ID)
	d.Set("identity_provider_id", group.IdentityProviderId)
	d.Set("name", group.Name)
	return resourceReadGroup(ctx, d, m)
}

func resourceReadGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	groupID := d.Id()
	group, err := client.Groups.GetGroup(groupID, IDPID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to get group"))
	}
	d.SetId(group.ID)
	d.Set("identity_provider_id", group.IdentityProviderId)
	d.Set("name", group.Name)
	return nil
}

func resourceDeleteGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting group")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}

	groupID := d.Id()
	err := client.Groups.DeleteGroup(IDPID, groupID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to delete group"))
	}
	d.SetId("")
	log.Printf("[INFO] Done Deleting group")
	return nil
}
