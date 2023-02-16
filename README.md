# Terraform Provider Loft (terraform-provider-loft)
Manage Loft Spaces and Virtual Clusters using terraform.

## Using the provider

You can browse the documentation on the [Terraform provider registry](https://registry.terraform.io/providers/loft-sh/loft/latest/docs)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 1.0.x
-	[Go](https://golang.org/doc/install) >= 1.17

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

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Generating Provider Code

In order to make it easier to add new resources to the provider, it is possible to generate much of the resource and data source implementations.

Because this works using the OpenAPI schema, information about whether a nested property is a pointer or struct is lost.

In addition, properties computed or automatically filled by controllers must be manually configured using the `Computed: true` terraform schema.

Finally, the code generation templates, located in the `gen/templates` directory, do not currently distinguish between cluster or namespace scoped resources.

Due to those exceptions the generated code must be corrected until the provider compiles successfully and tested against the local Loft deployment to ensure that the computed values can be read and written correctly. Below are the recommended steps for generating a new resource.

1. Open the file `GNUmakefile`, and update `generate-models` to include your desired model using the `--model` flag. Here's an example of updating the current command with the `Cluster` resource:
```sh
	swagger generate client \
		--spec ./gen/swagger.json \
		--template-dir ./gen/templates \
		--config-file ./gen/resources.yml \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.VirtualClusterInstance" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.SpaceInstance" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.Project" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.Cluster" \
		> gen/resources.log
```
2. Because this model may reference many other models, the command to generate the schemas must be updated to include the referenced models. The models can be determined by running `make generate-models` and `make build-local` and updating the command based on the missing models as reported in the compile errors. In this case, the following models need to be added:
```sh
        ...
        --model "com.github.loft-sh.api.v3.pkg.apis.management.v1.ClusterSpec" \
        --model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.SecretRef" \
        ...
```
3. Continue correcting compile errors until all required models have been generated. Note that the code generation will not update existing files after they are created.
4. In addition, the generated code may need to be updated for namespace vs cluster scoped resources. In this example, the `namespace` is not required since this is a cluster scoped resource.
5. Next, add a test in the `tests` directory to ensure that terraform can create, update, delete, import and read using the generated code. Some properties may need to be configured as `Computed` if the terraform test consistently shows that resources are our of sync and would be updated multiple times.