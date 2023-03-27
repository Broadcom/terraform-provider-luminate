package provider

import (
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
)

func LuminateCollectionSiteLink() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"links": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Site ID",
						},
						"collection_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Link ID",
						},
					},
				},
			},
		},
		Create: resourceCollectionSiteLinkCreate,
		Delete: resourceCollectionSiteLinkDelete,
		Read:   resourceCollectionSiteLinkRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCollectionSiteLinkRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Reading site link")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	links := *extractCollectionSiteLinkFields(d)
	if len(links) != 1 {
		errMessage := fmt.Sprintf("unable to get site link, wrong number of links: %d", len(links))
		return errors.New(errMessage)
	}
	_, err := client.CollectionAPI.GetCollectionSiteLinks(links[0].CollectionID)
	if err != nil {
		return err
	}
	d.SetId("site_link")
	return nil
}

func resourceCollectionSiteLinkCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating site link")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	links := extractCollectionSiteLinkFields(d)

	_, err := client.CollectionAPI.LinkSiteToCollection(*links)
	if err != nil {
		return err
	}
	return nil
}

func resourceCollectionSiteLinkDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Deleting site link")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	links := *extractCollectionSiteLinkFields(d)
	if len(links) != 1 {
		errMessage := fmt.Sprintf("unable to delete site link, wrong number of links: %d", len(links))
		return errors.New(errMessage)
	}
	err := client.CollectionAPI.UnlinkSiteFromCollection(links[0])
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func extractCollectionSiteLinkFields(d *schema.ResourceData) *[]dto.CollectionSiteLink {
	k, ok := d.Get("links").([]interface{})
	links := make([]dto.CollectionSiteLink, 0)
	if ok && len(k) > 0 {
		for _, v := range k {
			link := v.(map[string]interface{})
			links = append(links, dto.CollectionSiteLink{
				CollectionID: link["collection_id"].(string),
				SiteID:       link["site_id"].(string),
			})
		}
	}

	return &links
}
