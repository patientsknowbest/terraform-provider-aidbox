package provider

import (
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testProviderFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var (
	testProvider          *schema.Provider
	testProviderFactories map[string]func() (*schema.Provider, error)
)

func envOrDefault(env string, df string) string {
	e, isPresent := os.LookupEnv(env)
	if isPresent {
		return e
	} else {
		return df
	}
}

func init() {
	apiClient := aidbox.NewApiClient(envOrDefault("AIDBOX_URL", "http://localhost:8888"), envOrDefault("AIDBOX_CLIENT", "root"), envOrDefault("AIDBOX_CLIENT_SECRET", "secret"))
	testProvider = New(apiClient)()
	testProviderFactories = map[string]func() (*schema.Provider, error){
		"aidbox": func() (*schema.Provider, error) {
			return testProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := testProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

// Used to skip tests that require a legacy mode of operation / incompatible with schema mode
func requireLegacyMode(t *testing.T) {
	if os.Getenv("TF_ACC_AIDBOX_MODE") != "legacy" {
		t.Skip("Test is not compatible with schema mode, skipping")
	}
}

// Used to skip tests that are incompatible with older modes of operation / schema mode
func requireSchemaMode(t *testing.T) {
	if os.Getenv("TF_ACC_AIDBOX_MODE") == "legacy" {
		t.Skip("Test is not compatible with non-schema mode, skipping")
	}
}
