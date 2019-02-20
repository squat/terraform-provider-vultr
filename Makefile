.PHONY: all-release build clean errcheck fmt fmt-go fmt-terraform fmtcheck lint lint-go lint-terraform release test test-compile testacc sweep vendor vendor-status vet

ARCH ?= amd64
PLATFORM ?= linux
ALL_PLATFORMS := darwin linux windows
BIN := terraform-provider-vultr
PKG := github.com/squat/$(BIN)
BUILD_IMAGE ?= golang:1.10.0-alpine
TEST ?= $$(go list ./... | grep -v 'vendor/')
GOFMT_FILES ?= $$(find . -name '*.go' | grep -v vendor)
SRC := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
SWEEP?=us-east-1,us-west-2
TERRAFORMFMT_FILES ?= examples
TESTARGS ?=
TAG := $(shell git describe --abbrev=0 --tags HEAD 2>/dev/null)
COMMIT := $(shell git rev-parse HEAD)
VERSION := $(COMMIT)
ifneq ($(TAG),)
    ifeq ($(COMMIT), $(shell git rev-list -n1 $(TAG)))
        VERSION := $(TAG)
    endif
endif
DIRTY := $(shell test -z "$$(git diff --shortstat 2>/dev/null)" || echo -dirty)
VERSION := $(VERSION)$(DIRTY)

default: build

build: fmtcheck
	go install

all-release: $(addprefix release-, $(ALL_PLATFORMS))

release-%:
	@$(MAKE) --no-print-directory ARCH=$(ARCH) PLATFORM=$* release

release: bin/$(BIN)_$(VERSION)_$(PLATFORM)_$(ARCH).tar.gz.asc

bin/$(PLATFORM)/$(ARCH):
	@mkdir -p bin/$(PLATFORM)/$(ARCH)

bin/$(BIN)_$(VERSION)_$(PLATFORM)_$(ARCH).tar.gz.asc: bin/$(BIN)_$(VERSION)_$(PLATFORM)_$(ARCH).tar.gz
	@cd bin && gpg --armor --detach-sign $(<F)

bin/$(BIN)_$(VERSION)_$(PLATFORM)_$(ARCH).tar.gz: bin/$(PLATFORM)/$(ARCH)/$(BIN)_$(VERSION)
	@tar -czf $@ -C $(<D) $(<F)

bin/$(PLATFORM)/$(ARCH)/$(BIN)_$(VERSION): $(SRC) glide.yaml bin/$(PLATFORM)/$(ARCH)
	@echo "building: $@"
	@docker run --rm \
	    -u $$(id -u):$$(id -g) \
	    -v $$(pwd):/go/src/$(PKG) \
	    -v $$(pwd)/bin/$(PLATFORM)/$(ARCH):/go/bin \
	    -w /go/src/$(PKG) \
	    $(BUILD_IMAGE) \
	    /bin/sh -c " \
	        GOARCH=$(ARCH) \
	        GOOS=$(PLATFORM) \
		CGO_ENABLED=0 \
		go build -o /go/bin/$(BIN)_$(VERSION) \
	    "	

fmt: fmt-go fmt-terraform

fmt-go:
	gofmt -w -s $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt-terraform:
	terraform fmt $(TERRAFORMFMT_FILES)

lint: lint-go lint-terraform

lint-go:
	@echo 'golint $(TEST)'
	@lint_res=$$(golint $(TEST)); if [ -n "$$lint_res" ]; then \
		echo ""; \
		echo "Golint found style issues. Please check the reported issues"; \
		echo "and fix them if necessary before submitting the code for review:"; \
		echo "$$lint_res"; \
		exit 1; \
	fi
	@echo 'gofmt -d -s $(GOFMT_FILES)'
	@fmt_res=$$(gofmt -d -s $(GOFMT_FILES)); if [ -n "$$fmt_res" ]; then \
		echo ""; \
		echo "Gofmt found style issues. Please check the reported issues"; \
		echo "and fix them if necessary before submitting the code for review:"; \
		echo "$$fmt_res"; \
		exit 1; \
	fi

lint-terraform:
	@echo "terraform fmt --check=true $(TERRAFORMFMT_FILES)"
	@lint_res=$$(terraform fmt --check=true $(TERRAFORMFMT_FILES)); if [ -n "$$lint_res" ]; then \
		echo ""; \
		echo "Terraform fmt found style issues. Please check the reported issues"; \
		echo "and fix them if necessary before submitting the code for review:"; \
		echo "$$lint_res"; \
		exit 1; \
	fi

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test -i $(TEST) || exit 1
	go test $(TESTARGS) -timeout=30s -parallel=4 $(TEST)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo 'go vet $(TEST)'
	@go vet $(TEST); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor:
	@glide install -v
	@glide-vc --only-code --no-tests

vendor-status:
	@glide list

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

clean:
	@rm -rf bin
