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

Trial license can be obtained as per [aidbox documentation](https://docs.aidbox.app/getting-started/installation/setup-aidbox.dev)

```sh
$ (cd scripts && AIDBOX_LICENSE_ID=<your-devbox-license-id> AIDBOX_LICENSE_KEY=<your-devbox-license-key> MB_LICENSE=<your-multibox-license> docker-compose up -d)
$ make testacc
```
