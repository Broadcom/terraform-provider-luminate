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

func LuminateUser() *schema.Resource {
	return &schema.Resource{
		Schema:        CommonUserDataSchema(),
		ReadContext:   resourceReadUsers,
		DeleteContext: resourceDeleteUser,
		CreateContext: resourceReadOnlyDeleteCreate,
		UpdateContext: resourceReadOnlyDeleteCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
func resourceReadOnlyDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceReadUsers(ctx, d, meta) // Return the resource read function
}

func resourceDeleteUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting collection")
	client, ok := m.(*service.LuminateService)
	if !ok {
		return diag.FromErr(errors.New("invalid client"))
	}

	userIds := d.Get("user_ids").([]interface{})
	for _, userID := range userIds {
		userIDString := userID.(string)
		err := client.Users.DeleteUser(IDPID, userIDString)
		if err != nil {
			return diag.FromErr(errors.Wrap(err, fmt.Sprintf("failed to delete user with ID: %s", userIDString)))
		}
	}
	d.SetId("")
	log.Printf("[INFO] Done Deleting user")
	return nil
}
