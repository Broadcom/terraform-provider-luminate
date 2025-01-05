// Copyright (c) Symantec ZTNA
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateCollection() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection Name",
			},
		},
		CreateContext: resourceCreateCollection,
		ReadContext:   resourceReadCollection,
		UpdateContext: resourceUpdateCollection,
		DeleteContext: resourceDeleteCollection,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating collection")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client "))
	}
	collectionName := d.Get("name").(string)
	collection, err := client.CollectionAPI.CreateCollection(collectionName)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to create collection"))
	}
	d.SetId(collection.ID.String())
	return resourceReadCollection(ctx, d, m)
}

func resourceReadCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading colelction")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	collectionID := d.Id()
	collection, err := client.CollectionAPI.GetCollection(collectionID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to get collection"))
	}
	d.SetId(collection.ID.String())
	return nil
}

func resourceUpdateCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating collection")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	collectionName := d.Get("name").(string)
	collectionID := d.Id()
	_, err := client.CollectionAPI.UpdateCollection(collectionName, collectionID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to update collection"))
	}
	return resourceReadCollection(ctx, d, m)
}

func resourceDeleteCollection(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting collection")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}
	collectionID, err := uuid.Parse(d.Id())
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to parse collection id"))
	}
	err = client.CollectionAPI.DeleteCollection(collectionID.String())
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to delete collection"))
	}
	return nil
}
