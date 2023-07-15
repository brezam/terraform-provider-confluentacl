package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAclCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("0.15.4"))),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAclConfig(testRealResource.SaName, testRealResource.EnvId, testRealResource.ClusterId, testRealResource.RestEndpoint),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("confluentacl_confluentacl.example", "id"),
					resource.TestCheckResourceAttr("confluentacl_acl.example", "cluster_id", testRealResource.ClusterId),
					resource.TestCheckResourceAttr("confluentacl_acl.example", "service_account_name", testRealResource.SaName),
				),
			},
		},
	})
}

func testAccAclConfig(saName, envId, resourceId, restEndpoint string) string {
	return fmt.Sprintf(`
		resource "confluentacl_api_key" "example" {
				service_account_name = "%s"
				environment_id       = "%s"
				resource_id          = "%s"
			}

			resource "confluentacl_acl" "example" {
				rest_endpoint        = "%s"
				cluster_id           = confluentacl_api_key.example.resource_id
				service_account_name = confluentacl_api_key.example.service_account_name
			
				resource_type = "TOPIC"
				resource_name = "test"
				pattern_type  = "PREFIXED"
				host          = "*"
				operation     = "READ"
				permission    = "ALLOW"
			}
		`, saName, envId, resourceId, restEndpoint)
}
