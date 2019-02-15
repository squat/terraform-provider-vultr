package vultr

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProvider = Provider().(*schema.Provider)
var testAccProviders = map[string]terraform.ResourceProvider{
	"vultr": testAccProvider,
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("expected provider to validate: %v", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("VULTR_API_KEY") == "" {
		t.Fatal("VULTR_API_KEY must be set for acceptance tests")
	}
}
