GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

default: testacc

.PHONY: build-local
build-local:
	go build -o "${HOME}/.terraform.d/plugins/registry.terraform.io/loft-sh/loft/0.0.1/$(GOOS)_$(GOARCH)/terraform-provider-loft_v0.0.1"

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 TF_ACC_PROVIDER_NAMESPACE='loft-sh' go test ./... -v $(TESTARGS) -timeout 120m
