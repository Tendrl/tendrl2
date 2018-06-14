# VARIABLES
PACKAGE="github.com/Tendrl/tendrl2"
BINARY_NAME="tendrl2"

default: usage

clean: ## Trash binary files
	@echo "--> cleaning..."
	@go clean || (echo "Unable to clean project" && exit 1)
	@rm -rf $(GOPATH)/bin/$(BINARY_NAME) 2> /dev/null
	@curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	@dep ensure
	@echo "Clean OK"

test: ## Run all tests
	@echo "--> testing..."
	@go test -v $(PACKAGE)/...

install: clean ## Compile sources and build binary and install binary
	@echo "--> installing..."
	@go install $(PACKAGE) || (echo "Compilation error" && exit 1)
	@echo "Install OK"

run: install ## Run your application
	@echo "--> running application..."
	@$(GOPATH)/bin/$(BINARY_NAME)

usage: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: clean ## Compile sources and build binary
	@echo "--> building..."
	@go build $(PACKAGE) || (echo "Compilation error" && exit 1)
	@echo "Build OK"
