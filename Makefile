.PHONY: build fmt lint test vendor vendor-status vet 

TEST?=$$(go list ./... | grep -v 'vendor/')
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: build

build:
	go install

fmt:
	gofmt -w -s $(GOFMT_FILES)

lint:
	@echo "golint ."
	@lint_res=$$(go list ./... | grep -v vendor/ | xargs -n 1 golint) ; if [ -n "$$lint_res" ]; then \
		echo ""; \
		echo "Golint found style issues. Please check the reported issues"; \
		echo "and fix them if necessary before submitting the code for review:"; \
		echo "$$lint_res"; \
		exit 1; \
	fi

test: vet lint
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

vendor:
	@dep ensure

vendor-status:
	@dep status

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi