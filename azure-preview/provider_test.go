package azurepreview

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"azurepreview": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	/*if err := os.Getenv("AZURE_CLIENT_ID"); err == "" {
		t.Fatal("AZURE_CLIENT_ID must be set for acceptance tests")
	}

	if err := os.Getenv("AZURE_CLIENT_SECRET"); err == "" {
		t.Fatal("AZURE_CLIENT_ID must be set for acceptance tests")
	}

	if err := os.Getenv("AZURE_TENANT_ID"); err == "" {
		t.Fatal("AZURE_TENANT_ID must be set for acceptance tests")
	}*/

	if err := os.Getenv("AZURE_TEST_ENROLLMENT_ACCOUNT"); err == "" {
		t.Fatal("AZURE_TEST_ENROLLMENT_ACCOUNT must be set for acceptance tests")
	}
}
