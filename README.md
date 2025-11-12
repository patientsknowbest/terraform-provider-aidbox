# Terraform Provider Aidbox

Terraform provider for [aidbox](https://docs.aidbox.app/).

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

See [examples](examples/) directory.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

Acceptance tests require an aidbox server to run against. 

You can start aidbox in docker with the provided [docker-compose](scripts/docker-compose.yaml) file.

Trial license can be obtained either
- as per [aidbox documentation](https://docs.aidbox.app/overview/aidbox-user-portal/licenses)
- we can also issue our own development licenses from now on, ask around for these, or try accessing the [user portal](https://aidbox.app/ui/portal#/project/f07750f6-28e3-44be-a8f8-c2004ef2b1ea/license)

```sh
(cd scripts && AIDBOX_LICENSE=<your-aidbox-license> docker-compose up -d)
make testacc

# if you have access to the pkb 1pw vault
(cd scripts && AIDBOX_LICENSE=$(op read op://PHR-dev-e2e/aidbox-dev-license/license-key) docker-compose up -d)
make testacc
```

Erase the test cache in case you don't have code changes but still want to run the tests from scratch
(e.g. you only have environmental changes, like updating the docker image):
`go clean -testcache`

## Testing the provider

### Testing the provider with legacy modes of operation

In case you want to test functionality that predates schema mode and is not compatible with it, use the compose
file docker-compose-legacy.yaml
```shell
docker compose -f scripts/docker-compose.yaml up
TF_ACC_AIDBOX_MODE=legacy make testacc
```

### Testing the provider with observing HTTP requests

Often you can get yourself into states you don't understand how you go into. To validate what's exactly happening
under the hood to validate your assumptions the HTTP requests-responses the provider code and terraform's test
framework is making can be dumped if the env `TF_ACC_DUMP_HTTP` is set to true

### Trying out the provider without releasing

Once the provider is in a suitable state, further to the above testing you can try it out by pointing terraform to look
for your build of the provider and use that instead of a released version. Put this in a CLI configuration file (e.g.
into a file called .terraformrc in your home dir)

```terraform
# https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers
# override all provider installations
provider_installation {
  # override this specific provider
  dev_overrides {
    "patientsknowbest/aidbox" = "/home/plugin-developer/go/bin"
  }
  # must also tell terraform how to look for other providers
  direct {}
}
```
