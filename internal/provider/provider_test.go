package provider

import (
	"github.com/patientsknowbest/terraform-provider-aidbox/internal/aidbox"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testProviderFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var (
	testProvider                  *schema.Provider
	testMultiboxProvider          *schema.Provider
	testProviderFactories         map[string]func() (*schema.Provider, error)
	testMultiboxProviderFactories map[string]func() (*schema.Provider, error)
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
	client := aidbox.NewClient(envOrDefault("AIDBOX_URL", "http://localhost:8888"), envOrDefault("AIDBOX_CLIENT", "root"), envOrDefault("AIDBOX_CLIENT_SECRET", "secret"), false)
	testProvider = New(client)()
	testProviderFactories = map[string]func() (*schema.Provider, error){
		"aidbox": func() (*schema.Provider, error) {
			return testProvider, nil
		},
	}

	client2 := aidbox.NewClient(envOrDefault("MULTIBOX_URL", "http://localhost:8889"), envOrDefault("AIDBOX_CLIENT", "root"), envOrDefault("AIDBOX_CLIENT_SECRET", "secret"), true)
	testMultiboxProvider = New(client2)()
	testMultiboxProviderFactories = map[string]func() (*schema.Provider, error){
		"aidbox": func() (*schema.Provider, error) {
			return testMultiboxProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := testProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMultiboxProvider(t *testing.T) {
	if err := testMultiboxProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
