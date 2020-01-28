package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

const testAccSshGwApplication = `
resource "luminate_site" "new-site" {
  name = "tf-site-upd2"
}

data "luminate_aws_integration" "my-aws_integration" {
	  integration_name = "terraform-acceptance"
}

resource "luminate_ssh_gw_application" "new-ssh-gw-application" {
  site_id = "${luminate_site.new-site.id}"
  name = "tf-ssh-gw"
  
  integration_id = "${data.luminate_aws_integration.my-aws_integration.id}"

  tags = {
    Type = "prod"
	Name = "my-name"
  }

  vpc {
    region = "eu-west-1"
    cidr_block = "127.0.0.1/14"
    vpc_id = "pc-01229e075c14c11a9"
  }

  vpc {
    region = "us-west-1"
    cidr_block = "1.1.1.1/14"
    vpc_id = "pc-01229e075c14c11a9"
  }
}
`

const testAccSshGwApplicationUpdate = `
resource "luminate_site" "new-site" {
  name = "tf-site-upd2"
}

data "luminate_aws_integration" "my-aws_integration" {
	  integration_name = "terraform-acceptance"
}

resource "luminate_ssh_gw_application" "new-ssh-gw-application" {
  site_id = "${luminate_site.new-site.id}"
  name = "tf-ssh-gw"
  
  integration_id = "${data.luminate_aws_integration.my-aws_integration.id}"

  tags = {
    Type = "prod-2"
	Name = "my-name"
  }

  vpc {
    region = "eu-west-3"
    cidr_block = "127.0.0.1/14"
    vpc_id = "pc-01229e075c14c11a9"
  }

  vpc {
    region = "us-west-2"
    cidr_block = "1.1.1.1/14"
    vpc_id = "pc-01229e075c14c11a9"
  }
}
`

func TestAccLuminateSshGwApplication(t *testing.T) {
	resourceName := "luminate_ssh_gw_application.new-ssh-gw-application"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccSleep,
				Config:    testAccSshGwApplication,
				Destroy:   false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-ssh-gw"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tf-ssh-gw.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tf-ssh-gw.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "tags.Type", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "my-name"),
					resource.TestCheckResourceAttr(resourceName, "vpc.0.region", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceName, "vpc.0.cidr_block", "127.0.0.1/14"),
					resource.TestCheckResourceAttr(resourceName, "vpc.0.vpc_id", "pc-01229e075c14c11a9"),
					resource.TestCheckResourceAttr(resourceName, "vpc.1.region", "us-west-1"),
					resource.TestCheckResourceAttr(resourceName, "vpc.1.cidr_block", "1.1.1.1/14"),
					resource.TestCheckResourceAttr(resourceName, "vpc.1.vpc_id", "pc-01229e075c14c11a9"),
				),
			},
			{
				PreConfig: testAccSleep,
				Config:    testAccSshGwApplicationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-ssh-gw"),
					resource.TestCheckResourceAttr(resourceName, "visible", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_address", fmt.Sprintf("tf-ssh-gw.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "luminate_address", fmt.Sprintf("tf-ssh-gw.%s", testAccDomain)),
					resource.TestCheckResourceAttr(resourceName, "tags.Type", "prod-2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "my-name"),
					resource.TestCheckResourceAttr(resourceName, "vpc.0.region", "eu-west-3"),
					resource.TestCheckResourceAttr(resourceName, "vpc.0.cidr_block", "127.0.0.1/14"),
					resource.TestCheckResourceAttr(resourceName, "vpc.0.vpc_id", "pc-01229e075c14c11a9"),
					resource.TestCheckResourceAttr(resourceName, "vpc.1.region", "us-west-2"),
					resource.TestCheckResourceAttr(resourceName, "vpc.1.cidr_block", "1.1.1.1/14"),
					resource.TestCheckResourceAttr(resourceName, "vpc.1.vpc_id", "pc-01229e075c14c11a9"),
				),
			},
		},
	})
}
