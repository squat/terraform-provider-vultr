.PHONY: build fmt fmt-go fmt-terraform lint lint-go lint-terraform test vendor vendor-status vet 

TEST?=$$(go list ./... | grep -v 'vendor/')
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
TERRAFORMFMT_FILES?=examples
TESTARGS?=

default: build

build:
	go install

fmt: fmt-go fmt-terraform

fmt-go:
	gofmt -w -s $(GOFMT_FILES)

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
	@echo "terraform fmt $(TERRAFORMFMT_FILES)"
	@lint_res=$$(terraform fmt $(TERRAFORMFMT_FILES)); if [ -n "$$lint_res" ]; then \
		echo ""; \
		echo "Terraform fmt found style issues. Please check the reported issues"; \
		echo "and fix them if necessary before submitting the code for review:"; \
		echo "$$lint_res"; \
		exit 1; \
	fi


test: vet lint
	go test -i $(TEST) || exit 1
	go test $(TESTARGS) -timeout=30s -parallel=4 $(TEST)

vendor:
	@glide install -v

vendor-status:
	@glide list

vet:
	@echo 'go vet $(TEST)'
	@go vet $(TEST); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
