donothing:
	@echo "Usage:"
	@echo "make build"
	@echo "make run"

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

build: guard-GOPATH
	mkdir -p $$GOPATH/bin/linux
	mkdir -p $$GOPATH/bin/darwin
	GOOS=linux GOARCH=amd64 go build -v -o $$GOPATH/bin/linux/go-cli $$GOPATH/src/github.com/nmls/go-cli/*.go
	GOOS=darwin GOARCH=amd64 go build -v -o $$GOPATH/bin/darwin/go-cli $$GOPATH/src/github.com/nmls/go-cli/*.go

run_linux: guard-GOPATH
	@$$GOPATH/bin/linux/go-cli

run_darwin: guard-GOPATH
	@$$GOPATH/bin/darwin/go-cli
