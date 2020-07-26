package azurepreview

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAzurePreviewSubscription_basic(t *testing.T) {
	name := fmt.Sprintf("testacc-%s", acctest.RandString(6))
	enrollmentAccount := os.Getenv("AZURE_TEST_ENROLLMENT_ACCOUNT")
	offerType := "MS-AZR-0017P"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzurePreviewSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAzurePreviewSubscriptionConfigBasic(name, enrollmentAccount, offerType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzurePreviewSubscriptionExists("azurepreview_subscription.test"),
				),
			},
		},
	})
}

func testAccCheckAzurePreviewSubscriptionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Meta).Subscriptions
	ctx := testAccProvider.Meta().(*Meta).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurepreview_subscription" {
			continue
		}

		subscriptionID, err := parseSubscriptionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, subscriptionID)
		if err != nil {
			if resp.IsHTTPStatus(404) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Subscription ID still exists: %s", *resp.SubscriptionID)
	}

	return nil
}

func testAccCheckAzurePreviewSubscriptionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subscription ID set")
		}

		client := testAccProvider.Meta().(*Meta).Subscriptions
		ctx := testAccProvider.Meta().(*Meta).StopContext

		subscriptionID, err := parseSubscriptionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.Get(ctx, subscriptionID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAzurePreviewSubscriptionConfigBasic(name, enrollmentAccount, offerType string) string {
	return fmt.Sprintf(`
resource "azurepreview_subscription" "test" {
  name               = "%s"
  enrollment_account = "%s"
  offer_type         = "%s"
}
`, name, enrollmentAccount, offerType)
}
