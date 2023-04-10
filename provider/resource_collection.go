package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
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
		Create: resourceCreateCollection,
		Read:   resourceReadCollection,
		Update: resourceUpdateCollection,
		Delete: resourceDeleteCollection,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateCollection(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating collection")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	collectionName := d.Get("name").(string)
	collection, err := client.CollectionAPI.CreateCollection(collectionName)
	if err != nil {
		return errors.Wrap(err, "failed to create collection")
	}
	d.SetId(collection.ID.String())
	return nil
}

func resourceReadCollection(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Reading colelction")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	collectionID := d.Id()
	collection, err := client.CollectionAPI.GetCollection(collectionID)
	if err != nil {
		return errors.Wrap(err, "failed to get collection")
	}
	d.SetId(collection.ID.String())
	return nil
}

func resourceUpdateCollection(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Updating collection")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	collectionName := d.Get("name").(string)
	collectionID := d.Id()
	_, err := client.CollectionAPI.UpdateCollection(collectionName, collectionID)
	if err != nil {
		return errors.Wrap(err, "failed to update collection")
	}
	return nil
}

func resourceDeleteCollection(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Deleting collection")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	collectionID, err := uuid.Parse(d.Id())
	if err != nil {
		return errors.Wrap(err, "failed to parse collection id")
	}
	err = client.CollectionAPI.DeleteCollection(collectionID.String())
	if err != nil {
		return errors.Wrap(err, "failed to delete collection")
	}
	return nil
}
