SHELL=/bin/bash
.SHELLFLAGS=-euo pipefail -c

# Dependency Versions
# need https://github.com/kubernetes-sigs/controller-tools/pull/557
CONTROLLER_GEN_VERSION:=9cd8c2840e842554b3264634e3956e9f4d1c5497
KIND_VERSION:=v0.10.0
YQ_VERSION:=v4@v4.7.0
GOIMPORTS_VERSION:=v0.1.0
GOLANGCI_LINT_VERSION:=v1.39.0

# Build Flags
export CGO_ENABLED:=0
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
SHORT_SHA=$(shell git rev-parse --short HEAD)
VERSION?=$(shell echo ${BRANCH} | tr / -)-${SHORT_SHA}
BUILD_DATE=$(shell date +%s)
MODULE:=github.com/thetechnick/package-operator
LD_FLAGS=-X $(MODULE)/internal/version.Version=$(VERSION) \
			-X $(MODULE)/internal/version.Branch=$(BRANCH) \
			-X $(MODULE)/internal/version.Commit=$(SHORT_SHA) \
			-X $(MODULE)/internal/version.BuildDate=$(BUILD_DATE)

UNAME_OS:=$(shell uname -s)
UNAME_ARCH:=$(shell uname -m)

# PATH/Bin
DEPENDENCIES:=.cache/dependencies/$(UNAME_OS)/$(UNAME_ARCH)
export GOBIN?=$(abspath .cache/dependencies/bin)
export PATH:=$(GOBIN):$(PATH)

# Config
KIND_KUBECONFIG:=.cache/e2e/kubeconfig
export KUBECONFIG?=$(abspath $(KIND_KUBECONFIG))
export GOLANGCI_LINT_CACHE=$(abspath .cache/golangci-lint)
API_BASE:=packages.thetechnick.ninja
export SKIP_TEARDOWN?=

# Container
IMAGE_ORG?=quay.io/nschiede
package_OPERATOR_MANAGER_IMAGE?=$(IMAGE_ORG)/package-operator-manager:$(VERSION)

# -------
# Compile
# -------

all: \
	bin/linux_amd64/package-operator-manager

bin/linux_amd64/%: GOARGS = GOOS=linux GOARCH=amd64

bin/%: generate manifests FORCE
	$(eval COMPONENT=$(shell basename $*))
	@echo -e -n "compiling cmd/$(COMPONENT)...\n  "
	$(GOARGS) go build -ldflags "-w $(LD_FLAGS)" -o bin/$* cmd/$(COMPONENT)/main.go
	@echo

FORCE:

# prints the version as used by build commands.
version:
	@echo $(VERSION)
.PHONY: version

clean:
	@rm -rf bin .cache
.PHONY: clean

# ------------
# Dependencies
# ------------

# setup kind
KIND:=$(DEPENDENCIES)/kind/$(KIND_VERSION)
$(KIND):
	@echo "installing kind $(KIND_VERSION)..."
	$(eval KIND_TMP := $(shell mktemp -d))
	@(cd "$(KIND_TMP)" \
		&& go mod init tmp \
		&& go get "sigs.k8s.io/kind@$(KIND_VERSION)" \
	) 2>&1 | sed 's/^/  /'
	@rm -rf "$(KIND_TMP)" "$(dir $(KIND))" \
		&& mkdir -p "$(dir $(KIND))" \
		&& touch "$(KIND)" \
		&& echo

# setup controller-gen
CONTROLLER_GEN:=$(DEPENDENCIES)/controller-gen/$(CONTROLLER_GEN_VERSION)
$(CONTROLLER_GEN):
	@echo "installing controller-gen $(CONTROLLER_GEN_VERSION)..."
	$(eval CONTROLLER_GEN_TMP := $(shell mktemp -d))
	@(cd "$(CONTROLLER_GEN_TMP)" \
		&& go mod init tmp \
		&& go get "sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)" \
	) 2>&1 | sed 's/^/  /'
	@rm -rf "$(CONTROLLER_GEN_TMP)" "$(dir $(CONTROLLER_GEN))" \
		&& mkdir -p "$(dir $(CONTROLLER_GEN))" \
		&& touch "$(CONTROLLER_GEN)" \
		&& echo

# setup yq
YQ:=$(DEPENDENCIES)/yq/$(YQ_VERSION)
$(YQ):
	@echo "installing yq $(YQ_VERSION)..."
	$(eval YQ_TMP := $(shell mktemp -d))
	@(cd "$(YQ_TMP)" \
		&& go mod init tmp \
		&& go get "github.com/mikefarah/yq/$(YQ_VERSION)" \
	) 2>&1 | sed 's/^/  /'
	@rm -rf "$(YQ_TMP)" "$(dir $(YQ))" \
		&& mkdir -p "$(dir $(YQ))" \
		&& touch "$(YQ)" \
		&& echo

# setup goimports
GOIMPORTS:=$(DEPENDENCIES)/goimports/$(GOIMPORTS_VERSION)
$(GOIMPORTS):
	@echo "installing goimports $(GOIMPORTS_VERSION)..."
	$(eval GOIMPORTS_TMP := $(shell mktemp -d))
	@(cd "$(GOIMPORTS_TMP)" \
		&& go mod init tmp \
		&& go get "golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)" \
	) 2>&1 | sed 's/^/  /'
	@rm -rf "$(GOIMPORTS_TMP)" "$(dir $(GOIMPORTS))" \
		&& mkdir -p "$(dir $(GOIMPORTS))" \
		&& touch "$(GOIMPORTS)" \
		&& echo

# setup golangci-lint
GOLANGCI_LINT:=$(DEPENDENCIES)/golangci-lint/$(GOLANGCI_LINT_VERSION)
$(GOLANGCI_LINT):
	@echo "installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	$(eval GOLANGCI_LINT_TMP := $(shell mktemp -d))
	@(cd "$(GOLANGCI_LINT_TMP)" \
		&& go mod init tmp \
		&& go get "github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)" \
	) 2>&1 | sed 's/^/  /'
	@rm -rf "$(GOLANGCI_LINT_TMP)" "$(dir $(GOLANGCI_LINT))" \
		&& mkdir -p "$(dir $(GOLANGCI_LINT))" \
		&& touch "$(GOLANGCI_LINT)" \
		&& echo

# installs all project dependencies
dependencies: \
	$(KIND) \
	$(CONTROLLER_GEN) \
	$(YQ) \
	$(GOIMPORTS) \
	$(GOLANGCI_LINT)
.PHONY: dependencies

# ----------
# Deployment
# ----------

# Run against the configured Kubernetes cluster in ~/.kube/config or $KUBECONFIG
run-package-operator-manager: generate fmt vet manifests
	go run -ldflags "-w $(LD_FLAGS)" \
		./cmd/package-operator-manager/main.go \
			-metrics-addr="0"
.PHONY: run-package-operator-manager

# Run against the configured Kubernetes cluster in ~/.kube/config or $KUBECONFIG
run-coordination-operator-manager: generate fmt vet manifests
	go run -ldflags "-w $(LD_FLAGS)" \
		./cmd/coordination-operator-manager/main.go \
			-metrics-addr="0"
.PHONY: run-coordination-operator-manager

# ----------
# Generators
# ----------

# Generate manifests e.g. CRD, RBAC etc.
manifests: $(CONTROLLER_GEN)
	@echo "generating kubernetes manifests..."
	@controller-gen crd:crdVersions=v1,generateEmbeddedObjectMeta=true \
		rbac:roleName=package-operator-manager \
		paths="./..." \
		output:crd:artifacts:config=config/deploy 2>&1 | sed 's/^/  /'
	@echo

# Generate code
generate: $(CONTROLLER_GEN)
	@echo "generating code..."
	@controller-gen object paths=./apis/... 2>&1 | sed 's/^/  /'
	@echo

# Makes sandwich
# https://xkcd.com/149/
sandwich:
ifneq ($(shell id -u), 0)
	@echo "What? Make it yourself."
else
	@echo "Okay."
endif

# -------------------
# Testing and Linting
# -------------------

pre-commit-install: $(GOIMPORTS)
	@echo "installing pre-commit hooks using https://pre-commit.com/"
	@pre-commit install
.PHONY: pre-commit-install

fmt:
	go fmt ./...
.PHONY: fmt

vet:
	go vet ./...
.PHONY: vet

# Runs code-generators, checks for clean directory and lints the source code.
lint: generate fmt vet manifests $(GOLANGCI_LINT)
	@hack/validate-directory-clean.sh
	golangci-lint run ./... --deadline=15m
.PHONY: lint

# Runs unittests
test-unit: generate fmt vet manifests
	CGO_ENABLED=1 go test -race -v ./internal/... ./cmd/...
.PHONY: test-unit

# Runs the E2E testsuite against the currently selected cluster.
# FORCE_FLAGS ensures that the tests will not be cached
FORCE_FLAGS = -count=1
test-e2e: config/deploy/deployment.yaml
	@echo "running e2e tests..."

	@go test -v $(FORCE_FLAGS) ./e2e/...
.PHONY: test-e2e

# Sets up a local kind cluster and runs E2E tests against this local cluster.
test-e2e-local: export KUBECONFIG=$(abspath $(KIND_KUBECONFIG))
test-e2e-local: | setup-e2e-kind test-e2e
.PHONY: test-e2e-local

# Run the E2E testsuite after installing the packageOperator into the cluster.
test-e2e-ci: | apply-ao test-e2e

# make sure that we install our components into the kind cluster and disregard normal $KUBECONFIG
setup-e2e-kind: export KUBECONFIG=$(abspath $(KIND_KUBECONFIG))
setup-e2e-kind: | \
	create-kind-cluster \
	load-package-operator

create-kind-cluster: $(KIND)
	@echo "creating kind cluster package-operator-e2e..."
	@mkdir -p .cache/e2e
	@(source hack/determine-container-runtime.sh \
		&& $$KIND_COMMAND create cluster \
			--kubeconfig=$(KIND_KUBECONFIG) \
			--name="package-operator-e2e" \
		&& sudo chown $$USER: $(KIND_KUBECONFIG) \
		&& echo) 2>&1 | sed 's/^/  /'
.PHONY: create-kind-cluster

delete-kind-cluster: $(KIND)
	@echo "deleting kind cluster package-operator-e2e..."
	@(source hack/determine-container-runtime.sh \
		&& $$KIND_COMMAND delete cluster \
			--kubeconfig="$(KIND_KUBECONFIG)" \
			--name "package-operator-e2e" \
		&& rm -rf "$(KIND_KUBECONFIG)" \
		&& echo) 2>&1 | sed 's/^/  /'
.PHONY: delete-kind-cluster

# Load Package Operator Image into kind
load-package-operator: build-image-package-operator-manager
	@source hack/determine-container-runtime.sh \
		&& $$KIND_COMMAND load image-archive \
			.cache/image/package-operator-manager.tar \
			--name package-operator-e2e
.PHONY: load-package-operator

# Template deployment
config/deploy/deployment.yaml: FORCE $(YQ)
	@yq eval '.spec.template.spec.containers[0].image = "$(package_OPERATOR_MANAGER_IMAGE)"' \
			config/deploy/deployment.yaml.tpl > config/deploy/deployment.yaml

# Installs the package Operator into the kind e2e cluster.
apply-package-operator: $(YQ) load-package-operator config/deploy/deployment.yaml
	@echo "installing package Operator $(VERSION)..."
	@(source hack/determine-container-runtime.sh \
		&& kubectl apply -f config/deploy \
		&& echo -e "\nwaiting for deployment/package-operator..." \
		&& kubectl wait --for=condition=available deployment/package-operator -n package-operator --timeout=240s \
		&& echo) 2>&1 | sed 's/^/  /'
.PHONY: apply-package-operator

# ----------------
# Container Images
# ----------------

build-images: \
	build-image-package-operator-manager
.PHONY: build-images

push-images: \
	push-image-package-operator-manager
.PHONY: push-images

.SECONDEXPANSION:
build-image-%: bin/linux_amd64/$$*
	@echo "building image ${IMAGE_ORG}/$*:${VERSION}..."
	@(source hack/determine-container-runtime.sh \
		&& rm -rf ".cache/image/$*" ".cache/image/$*.tar" \
		&& mkdir -p ".cache/image/$*" \
		&& cp -a "bin/linux_amd64/$*" ".cache/image/$*" \
		&& cp -a "config/docker/$*.Dockerfile" ".cache/image/$*/Dockerfile" \
		&& cp -a "config/docker/passwd" ".cache/image/$*/passwd" \
		&& echo "building ${IMAGE_ORG}/$*:${VERSION}" \
		&& $$CONTAINER_COMMAND build -t "${IMAGE_ORG}/$*:${VERSION}" ".cache/image/$*" \
		&& $$CONTAINER_COMMAND image save -o ".cache/image/$*.tar" "${IMAGE_ORG}/$*:${VERSION}" \
		&& echo) 2>&1 | sed 's/^/  /'

push-image-%: build-image-$$*
	@echo "pushing image ${IMAGE_ORG}/$*:${VERSION}..."
	@(source hack/determine-container-runtime.sh \
		&& $$CONTAINER_COMMAND push "${IMAGE_ORG}/$*:${VERSION}" \
		&& echo pushed "${IMAGE_ORG}/$*:${VERSION}" \
		&& echo) 2>&1 | sed 's/^/  /'
