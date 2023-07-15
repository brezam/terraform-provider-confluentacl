package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestApiKeyCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("0.15.4"))),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApiKeyConfig(testRealResource.SaName, testRealResource.EnvId, testRealResource.ClusterId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("confluentacl_api_key.example", "id"),
				),
			},
		},
	})
}

func testAccApiKeyConfig(saName, envId, resourceId string) string {
	return fmt.Sprintf(`
		resource "confluentacl_api_key" "example" {
				service_account_name = "%s"
				environment_id       = "%s"
				resource_id          = "%s"
			}
		`, saName, envId, resourceId)
}
