package provider

import (
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"log"
	"sort"
)

func LuminateCollectionSiteLink() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID",
			},
			"collection_ids": {
				Type:     schema.TypeList,
				Required: true,
				StateFunc: func(v interface{}) string {
					list := v.([]string)
					sort.Strings(list)
					return fmt.Sprintf("%v", list)
				},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Create: resourceCollectionSiteLinkCreate,
		Delete: resourceCollectionSiteLinkDelete,
		Read:   resourceCollectionSiteLinkRead,
		Update: resourcesCollectionSiteLinkUpdate,
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
	siteID := d.Get("site_id").(string)

	res, err := client.CollectionAPI.GetCollectionsBySite(siteID)
	if err != nil {
		return err
	}

	d.SetId(siteID)
	ids := *res
	sort.Strings(ids)
	log.Println("collection_ids", ids)
	err = d.Set("collection_ids", ids)
	if err != nil {
		return errors.Wrapf(err, "unable to set collection_id for site %s", siteID)
	}
	return nil
}

func resourceCollectionSiteLinkCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating site link")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	links := extractCollectionSiteLinkFields(d)

	createdLinks, err := client.CollectionAPI.LinkSiteToCollection(*links)
	if err != nil || createdLinks == nil || *createdLinks == nil {
		return errors.Wrapf(err, "unable to link site to collections")
	}
	siteID := (*createdLinks)[0].SiteID
	collectionIDs := make([]string, len(*createdLinks))
	for i, link := range *createdLinks {
		collectionIDs[i] = link.CollectionID
	}

	d.SetId(siteID)

	sort.Strings(collectionIDs)
	log.Println("collection_ids", collectionIDs)
	err = d.Set("collection_ids", collectionIDs)
	if err != nil {
		return errors.Wrapf(err, "unable to set collection_id for site %s", siteID)
	}
	return nil
}

func resourcesCollectionSiteLinkUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Updating site link")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}
	if d.HasChange("collection_ids") {
		siteID := d.Get("site_id").(string)
		currentCollections, err := client.CollectionAPI.GetCollectionsBySite(siteID)
		if err != nil {
			return err
		}
		newCollectionState := d.Get("collection_ids").([]interface{})
		if newCollectionState == nil {
			return nil
		}
		// convert to string slice
		newCollectionStateStr := make([]string, len(newCollectionState))
		for i, c := range newCollectionState {
			newCollectionStateStr[i] = c.(string)
		}

		unlink, link := getUniqueValues(newCollectionStateStr, *currentCollections)
		if len(unlink) > 0 {
			for _, id := range unlink {
				err := client.CollectionAPI.UnlinkSiteFromCollection(dto.CollectionSiteLink{
					CollectionID: id,
					SiteID:       siteID,
				})
				if err != nil {
					return err
				}
			}
		}
		if len(link) > 0 {
			links := make([]dto.CollectionSiteLink, 0)
			for _, id := range link {
				links = append(links, dto.CollectionSiteLink{
					CollectionID: id,
					SiteID:       siteID,
				})
			}
			_, err := client.CollectionAPI.LinkSiteToCollection(links)
			if err != nil {
				return err
			}
		}
		sort.Strings(newCollectionStateStr)
		log.Println("collection_ids", newCollectionStateStr)
		err = d.Set("collection_ids", newCollectionStateStr)
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
	for _, link := range links {
		err := client.CollectionAPI.UnlinkSiteFromCollection(link)
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}

func extractCollectionSiteLinkFields(d *schema.ResourceData) *[]dto.CollectionSiteLink {
	siteID := d.Get("site_id").(string)
	collectionIDs, ok := d.Get("collection_ids").([]interface{})
	if !ok {
		return nil
	}
	links := make([]dto.CollectionSiteLink, 0)
	if len(collectionIDs) > 0 {
		for _, id := range collectionIDs {
			links = append(links, dto.CollectionSiteLink{
				CollectionID: id.(string),
				SiteID:       siteID,
			})
		}
	}

	return &links
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// getUniqueValues returns the values that are in the current state but not in the new state and vice versa
func getUniqueValues(newState []string, currentState []string) ([]string, []string) {
	var unlink []string
	var link []string
	for _, v := range currentState {
		if !contains(newState, v) {
			unlink = append(unlink, v)
		}
	}

	for _, v := range newState {
		if !contains(currentState, v) {
			link = append(link, v)
		}
	}

	return unlink, link
}
