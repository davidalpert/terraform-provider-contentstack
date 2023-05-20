<!-- PROJECT SHIELDS -->
<!--
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![License: GPL v3][license-shield]][license-url]
<!-- [![Issues][issues-shield]][issues-url] -->
<!-- [![Forks][forks-shield]][forks-url] -->
<!-- ![GitHub Contributors][contributors-shield] -->
<!-- ![GitHub Contributors Image][contributors-image-url] -->

<!-- PROJECT LOGO -->
<br />
<!-- vale Google.Headings = NO -->
<h1 align="center">terraform-plugin-contentstack-admin</h1>
<!-- vale Google.Headings = YES -->

<p align="center">
  A terraform plugin for managing ContentStack settings through the ContentStack Management API.
  <br />
  <a href="./README.md"><strong>README</strong></a>
  ·
  <a href="./CHANGELOG.md">CHANGELOG</a>
  <br />
  <!-- <a href="https://github.com/davidalpert/terraform-provider-contentstack-admin">View Demo</a>
  · -->
  <a href="https://github.com/davidalpert/terraform-provider-contentstack-admin/issues">Report Bug</a>
  ·
  <a href="https://github.com/davidalpert/terraform-provider-contentstack-admin/issues">Request Feature</a>
</p>

_This provider is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and is compatible with provider protocol v6._

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Install [task](https://taskfile.dev/install)

Then:

1. Clone the repository
1. Enter the repository directory
1. Build the provider
    ```shell
    task build
    ```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `task build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `task test-acceptance`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
task test-acceptance
```
