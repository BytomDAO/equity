PACKAGES    := $(shell go list ./... | grep -v '/vendor/')
BUILD_FLAGS := -ldflags "-X github.com/equity/compiler.GitCommit=`git rev-parse HEAD`"

all: test cmd equity

cmd:
	@echo "Building equitycmd to target/equitycmd"
	@go build $(BUILD_FLAGS) -o target/equitycmd compiler/cmd/equitycmd/equitycmd.go

equity:
	@echo "Building equity to target/equity"
	@go build $(BUILD_FLAGS) -o target/equity equity/main.go

clean:
	@rm -rf target/equitycmd
	@echo "Remove equitycmd successfully."
	@rm -rf target/equity
	@echo "Remove equity successfully."

test:
	@echo "====> Running go test"
	@go test -tags "equity" $(PACKAGES)

ci: test

.PHONY: all clean test ci cmd equity
