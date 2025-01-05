// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LuminateGroupUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "group id",
				ValidateFunc: utils.ValidateUuid,
				ForceNew:     true,
			},
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "user id",
				ValidateFunc: utils.ValidateUuid,
				ForceNew:     true,
			},
		},
		CreateContext: resourceCreateGroupUser,
		ReadContext:   resourceReadGroupUser,
		DeleteContext: resourceDeleteGroupUser,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateGroupUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE GROUP_USER CREATE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	if err := client.Groups.AssignUser(groupId, userId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(formatId(groupId, userId))

	return resourceReadGroupUser(ctx, d, m)
}

func resourceReadGroupUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE GROUP_USER READ")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)

	isAssigned, err := client.Groups.CheckAssignedUser(groupId, userId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !isAssigned {
		return diag.FromErr(errors.New("user wasn't assigned to group"))
	}

	return diagnostics
}

func resourceDeleteGroupUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE GROUP_USER DELETE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}

	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	if err := client.Groups.RemoveUser(groupId, userId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func formatId(groupId string, userId string) string {
	return fmt.Sprintf("%s_%s", groupId, userId)
}
