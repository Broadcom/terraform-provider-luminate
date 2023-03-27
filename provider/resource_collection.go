package provider

import (
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func LuminateCollection() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collection Name",
				Default:     false,
			},
		},
		Create: resourceCreateCollection,
		Read:   resourceReadCollection,
		Update: resourceUpdateCollection,
		Delete: resourceDeleteCollection,
	}
}

func resourceCreateCollection(d *schema.ResourceData, m interface{}) error {
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
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	collectionName := d.Get("name").(string)
	collection, err := client.CollectionAPI.GetCollection(collectionName)
	if err != nil {
		return errors.Wrap(err, "failed to get collection")
	}
	d.SetId(collection.ID.String())
	return nil
}

func resourceUpdateCollection(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	collectionName := d.Get("name").(string)
	collectionID := uuid.FromStringOrNil(d.Id())
	_, err := client.CollectionAPI.UpdateCollection(collectionID.String(), collectionName)
	if err != nil {
		return errors.Wrap(err, "failed to update collection")
	}
	return nil
}

func resourceDeleteCollection(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("invalid client")
	}
	collectionID := uuid.FromStringOrNil(d.Id())
	err := client.CollectionAPI.DeleteCollection(collectionID.String())
	if err != nil {
		return errors.Wrap(err, "failed to delete collection")
	}
	return nil
}
