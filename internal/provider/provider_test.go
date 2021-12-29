package provider

import (
	"github.com/hashicorp/terraform-provider-scaffolding/internal/aidbox"
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
	client := aidbox.NewClient(envOrDefault("AIDBOX_URL", "http://localhost:8888"), envOrDefault("AIDBOX_CLIENT", "root"), envOrDefault("AIDBOX_CLIENT_SECRET", "secret"))
	testProvider = New(client)()
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
