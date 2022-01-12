package provider

import (
	"errors"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service"
	"log"

	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

func LuminateGroupUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "group id",
				ValidateFunc: utils.ValidateUuid,
				ForceNew: true,
			},
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "user id",
				ValidateFunc: utils.ValidateUuid,
				ForceNew: true,
			},
		},
		Create: resourceCreateGroupsUser,
		Read:   resourceReadGroupsUser,
		Delete: resourceDeleteGroupsUser,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateGroupsUser(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE GROUP_USER CREATE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	if err := client.Groups.AssignUser(groupId, userId); err != nil {
		return err
	}

	d.SetId(formatId(groupId, userId))

	return nil
}

func resourceReadGroupsUser(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE GROUP_USER READ")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)

	isAssigned, err := client.Groups.CheckAssignedUser(groupId, userId)
	if err != nil {
		return err
	}

	if !isAssigned {
		return errors.New("user wasn't assigned to group")
	}

	return nil
}

func resourceDeleteGroupsUser(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LUMINATE GROUP_USER DELETE")

	client, ok := m.(*service.LuminateService)
	if !ok {
		return errors.New("unable to cast Luminate service")
	}

	groupId := d.Get("group_id").(string)
	userId := d.Get("user_id").(string)
	if err := client.Groups.RemoveUser(groupId, userId); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func formatId(groupId string, userId string) string {
	return fmt.Sprintf("%s_%s", groupId, userId)
}
