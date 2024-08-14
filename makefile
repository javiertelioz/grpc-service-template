# Export ==============================================================================================================
export LOCAL_BIN := $(CURDIR)/bin
export PATH := $(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help clean generate specs

help: ## ℹ️  Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Targets ==============================================================================================================

clean: ## 🧹 Remove generated files in proto directories
	@rm -f proto/helloworld/v1/*.{go,java,py,ts,yaml} || true
	@rm -f proto/payments/v1/*.{go,java,py,ts,yaml} || true
	@echo "🧹 Success clean folder"

generate: ## ⚙️ Generate code from proto files using buf
	@buf generate proto/helloworld/v1/*.proto
	@buf generate proto/payments/v1/*.proto
	@echo "⚙️ Success generate"

specs: ## 📜 Convert Swagger specs to OpenAPI v3 using swagger-cli
	@npx swagger-cli bundle proto/helloworld/v1/helloworld.swagger.yaml \
		--outfile proto/helloworld/v1/helloworld.openapi-v3.yaml --type yaml
	@npx swagger-cli bundle proto/payments/v1/payments.swagger.yaml \
		--outfile proto/payments/v1/payments.openapi-v3.yaml --type yaml
	@echo "📜 Success specs"

install: ## 🚀 Install tool and dependencies
	@cd proto && buf dep update && cd ../
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	@npm install -g @bufbuild/buf @bufbuild/protoc-gen-es @connectrpc/protoc-gen-connect-es
	@echo "🚀 Success Install"

linter: ## 🔍 Run golangci-lint on the Go codebase
	@golangci-lint run --config .code_quality/.golangci.yml
	@echo "🔍 Linting completed successfully"

commit: ## 📝 Create a commit using Commitizen
	@cz commit
	@echo "📝 Commit created successfully with Commitizen"
