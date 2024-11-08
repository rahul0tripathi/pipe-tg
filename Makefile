.EXPORT_ALL_VARIABLES:

include .env

GO := go
GOBIN	:= $(PWD)/_bin


.PHONY: build-common
build-common: ## - execute build common tasks clean and mod tidy
	@ $(GO) version
	@ $(GO) clean
	@ $(GO) mod tidy && $(GO) mod download
	@ $(GO) mod verify

build: build-common ## - build a debug binary to the current platform (windows, linux or darwin(mac))
	@ echo cleaning...
	@ rm -f $(GOBIN)/debug/$(OS)/$(service_name)
	@ echo building...
	@ $(GO) build -tags dev -o "$(GOBIN)/debug/$(OS)/$(service_name)" cmd/main.go
	@ ls -lah $(GOBIN)/debug/$(OS)/$(service_name)


.PHONY: build-linux-release
build-linux-release: build-common ## - build a static release linux elf(binary)
	@ CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags $(tags) -ldflags='-w -s -extldflags "-static"' -a -o "_bin/release/linux/$(service_name)" cmd/main.go

.PHONY: build-local-release
build-local-release: build-common ## - build a release binary to the current platform (windows, linux or darwin(mac))
	@ echo building
	@ @ CGO_ENABLED=0 go build -tags $(tags) -ldflags='-w -s -extldflags "-static"' -a -o "_bin/release/$(OS)/$(service_name)" cmd/main.go
	@ echo "_bin/release/$(OS)/"
	@ ls -lah _bin/release/$(OS)/$(service_name)
	@ echo "done"

.PHONY: test
test: build-common ## - execute go test command
	@ go test -v -cover ./...