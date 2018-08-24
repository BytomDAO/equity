PACKAGES    := $(shell go list ./... | grep -v '/vendor/')

all: test equitycmd

equitycmd:
	@echo "Building equitycmd to compiler/cmd/equitycmd/equitycmd"
	@go build -o compiler/cmd/equitycmd/equitycmd compiler/cmd/equitycmd/equitycmd.go

clean:
	@echo "Cleaning binaries built..."
	@rm -rf compiler/cmd/equitycmd/equitycmd
	@echo "Done."

test:
	@echo "====> Running go test"
	@go test -tags "equity" $(PACKAGES)

ci: test

.PHONY: all clean test ci
