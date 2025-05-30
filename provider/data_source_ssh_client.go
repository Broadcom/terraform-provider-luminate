// Copyright (c) Broadcom Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"
)

func LuminateDataSourceSshClient() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "ID is the id for this ssh-client",
				Computed:    true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "ssh-client to retrieve",
			},
			"description": {
				Type:        schema.TypeString,
				Description: "ssh-client's description",
				Computed:    true,
			},
			"key_size": {
				Type:        schema.TypeInt,
				Description: "The size of the key for the ssh-client",
				Computed:    true,
			},
			"created_on": {
				Type:        schema.TypeString,
				Description: "creation date of the ssh-client",
				Computed:    true,
			},
			"modified_on": {
				Type:        schema.TypeString,
				Description: "last time the ssh-client was modified",
				Computed:    true,
			},
			"last_accessed": {
				Type:        schema.TypeString,
				Description: "last time the ssh-client logged in",
				Computed:    true,
			},
			"expires": {
				Type:        schema.TypeString,
				Description: "the time the ssh-client key expired",
				Computed:    true,
			},
		},
		ReadContext: resourceReadSshClient,
	}
}

func resourceReadSshClient(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LUMINATE GET SSH-CLIENT")
	var diagnostics diag.Diagnostics
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("unable to cast Luminate service"))
	}
	sshClientName := d.Get("name").(string)

	sshClient, err := client.SshClientApi.GetSshClientByName(sshClientName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", sshClient.Description)
	d.Set("key_size", sshClient.KeySize)
	d.Set("created_on", sshClient.CreatedOn.Format(time.RFC3339))
	d.Set("modified_on", sshClient.ModifiedOn.Format(time.RFC3339))
	d.Set("last_accessed", sshClient.LastAccessed.Format(time.RFC3339))
	d.Set("expires", sshClient.Expires.Format(time.RFC3339))
	d.SetId(sshClient.Id)

	return diagnostics
}
