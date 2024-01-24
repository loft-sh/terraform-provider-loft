GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

default: testacc

.PHONY: generate-docs
generate-docs:
	go generate ./...

.PHONY: generate
generate-models:
	swagger generate client \
		--spec ./gen/swagger.json \
		--template-dir ./gen/templates \
		--config-file ./gen/resources.yml \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.VirtualClusterInstance" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.VirtualClusterTemplate" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.SpaceInstance" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.Project" \
		> gen/resources.log

	swagger generate client \
		--spec ./gen/swagger.json \
		--template-dir ./gen/templates \
		--config-file ./gen/schemas.yml \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.ProjectSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.SpaceInstanceSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.VirtualClusterInstanceSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.management.v1.VirtualClusterTemplateSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.Access" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.AllowedCluster" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.AllowedRunner" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.AllowedTemplate" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.AppParameter" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ArgoIntegrationSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ArgoProjectPolicyRule" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ArgoProjectRole" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ArgoProjectSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ArgoProjectSpecMetadata" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ArgoSSOSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.AutomaticImport" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.AutomaticImportVirtualClusters" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ClusterRef" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.ImportVirtualClustersSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.Member" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.NamespacePattern" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.Quotas" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.RancherIntegrationSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.RancherProjectRef" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.RequireTemplate" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.SpaceInstanceTemplateDefinition" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.SpaceTemplateDefinition" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.SyncMembersSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.TemplateMetadata" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.TemplateRef" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.UserOrTeam" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VaultIntegrationSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VaultAuthSpec" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VirtualClusterClusterRef" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VirtualClusterInstanceTemplateDefinition" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VirtualClusterTemplateDefinition" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VirtualClusterSpaceTemplateDefinition" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VirtualClusterTemplateSpaceTemplateRef" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.VirtualClusterTemplateVersion" \
		--model "com.github.loft-sh.api.v3.pkg.apis.storage.v1.WorkloadVirtualClusterTemplateDefinition" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.AppReference" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.ChartSecretRef" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.InstanceAccess" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.InstanceAccessRule" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.ProjectSecretRef" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.TemplateHelmChart" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.VirtualClusterAccessPoint" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.VirtualClusterAccessPointIngressSpec" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.VirtualClusterHelmChart" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.VirtualClusterHelmRelease" \
		--model "com.github.loft-sh.agentapi.v3.pkg.apis.loft.storage.v1.VirtualClusterProSpec" \
		--model "io.k8s.api.core.v1.SecretKeySelector" \
		> gen/schemas.log

.PHONY: build-local
build-local:
	go build -o "${HOME}/.terraform.d/plugins/registry.terraform.io/loft-sh/loft/0.0.1/$(GOOS)_$(GOARCH)/terraform-provider-loft_v0.0.1"

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 TF_ACC_PROVIDER_NAMESPACE='loft-sh' go test ./tests/... -v $(TESTARGS) -timeout 120m
