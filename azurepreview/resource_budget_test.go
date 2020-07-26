package azurepreview

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAzurePreviewBudget_basic(t *testing.T) {
	scope := fmt.Sprintf("subscriptions/%s", os.Getenv("AZURE_SUBSCRIPTION_ID"))
	name := fmt.Sprintf("testacc-%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzurePreviewBudgetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAzurePreviewBudgetConfigBasic(scope, name, acctest.RandString(6)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzurePreviewBudgetExists("azurepreview_budget.test"),
				),
			},
		},
	})
}

func testAccCheckAzurePreviewBudgetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Meta).Budgets
	ctx := testAccProvider.Meta().(*Meta).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurepreview_budget" {
			continue
		}

		id, err := parseBudgetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.Scope, id.BudgetName)
		if err != nil {
			if resp.IsHTTPStatus(404) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Budget ID still exists: %s", *resp.ID)
	}

	return nil
}

func testAccCheckAzurePreviewBudgetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Budget ID set")
		}

		client := testAccProvider.Meta().(*Meta).Budgets
		ctx := testAccProvider.Meta().(*Meta).StopContext

		id, err := parseBudgetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = client.Get(ctx, id.Scope, id.BudgetName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAzurePreviewBudgetConfigBasic(scope, name, random string) string {
	return fmt.Sprintf(`
resource "azurepreview_budget" "test" {
  scope      = "%s"
  name       = "%s"
  category   = "Cost"
  amount     = 1000
  time_grain = "BillingMonth"

  time_period {
    start_date = "2017-06-01T00:00:00Z"
    end_date   = "2035-06-01T00:00:00Z"
  }

  filters {
    resources = ["%s"]
    tag {
      name = "%s"
      values = [
        "%s"
      ]
    }
  }

  notification {
    name      = "%s"
    operator  = "GreaterThan"
    threshold = 80
    contact_roles = [
      "Contributor",
    ]
  }
}
`, scope, name, random, random, random, random)
}
