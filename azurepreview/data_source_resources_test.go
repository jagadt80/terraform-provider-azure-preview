package azurepreview

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAzurePreviewResources_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAzurePreviewResourcesConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azurepreview_resources.test", "resources.0.type", "Microsoft.Network/virtualNetworks"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAzurePreviewResourcesConfigBasic() string {
	return `
data "azurepreview_resources" "test" {
  type = "Microsoft.Network/virtualNetworks"
  tags = {
	is_spoke = true
  }
}
`
}
