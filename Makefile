ifndef GOOS
	GOOS := linux
endif

PACKAGES    := $(shell go list ./... | grep -v '/vendor/')
BUILD_FLAGS := -ldflags "-X github.com/equity/compiler.GitCommit=`git rev-parse HEAD`"
VERSION := $(shell awk -F= '/VersionMajor =/ {print $$2}' compiler/version.go | tr -d "\" ").$(shell awk -F= '/VersionMinor =/ {print $$2}' compiler/version.go | tr -d "\" ").$(shell awk -F= '/VersionPatch =/ {print $$2}' compiler/version.go | tr -d "\" ")
EQUITY_RELEASE := equity-$(VERSION)-$(GOOS)

all: test cmd equity

cmd:
	@echo "Building equitycmd to target/equitycmd"
	@go build $(BUILD_FLAGS) -o target/equitycmd compiler/cmd/equitycmd/equitycmd.go

equity:
	@echo "Building equity to target/equity"
	@go build $(BUILD_FLAGS) -o target/equity equity/main.go

ifeq ($(GOOS),windows)
release: equity
	cd target && cp -f equity $(EQUITY_RELEASE).exe
	cd target && md5sum $(EQUITY_RELEASE).exe > $(EQUITY_RELEASE).md5
	cd target && zip $(EQUITY_RELEASE).zip $(EQUITY_RELEASE).exe $(EQUITY_RELEASE).md5
	cd target && rm -f equity $(EQUITY_RELEASE).exe $(EQUITY_RELEASE).md5
else
release: equity
	cd target && cp -f equity $(EQUITY_RELEASE)
	cd target && md5sum $(EQUITY_RELEASE) > $(EQUITY_RELEASE).md5
	cd target && tar -czf $(EQUITY_RELEASE).tgz $(EQUITY_RELEASE) $(EQUITY_RELEASE).md5
	cd target && rm -f equity $(EQUITY_RELEASE) $(EQUITY_RELEASE).md5
endif

release-all: clean
	GOOS=linux   make release
	GOOS=windows make release

clean:
	@rm -rf target
	@echo "Cleaning target binaries successfully."

test:
	@echo "====> Running go test"
	@go test -tags "equity" $(PACKAGES)

ci: test

.PHONY: all clean test ci cmd equity
