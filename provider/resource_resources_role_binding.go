package provider

import "github.com/hashicorp/terraform/helper/schema"

func LuminateResourcesRoleBindings() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role",
			},
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Entity ID",
			},
			"subject_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subject Type",
			},
		},
		Create: resourceCreateRoleBinding,
		Read:   nil,
		Delete: nil,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
