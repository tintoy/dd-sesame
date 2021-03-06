default: version test build

fmt:
	go fmt github.com/tintoy/dd-sesame

# Peform a development (current-platform-only) build.
dev: fmt
	go build -o _bin/dd-sesame

install: dev
	go install

# Perform a full (all-platforms) build.
build: version build-windows64 build-linux64 build-mac64

build-windows64:
	GOOS=windows GOARCH=amd64 go build -o _bin/windows-amd64/dd-sesame.exe

build-linux64:
	GOOS=linux GOARCH=amd64 go build -o _bin/linux-amd64/dd-sesame

build-mac64:
	GOOS=darwin GOARCH=amd64 go build -o _bin/darwin-amd64/dd-sesame

# Produce archives for a GitHub release.
dist: build
	zip -9 _bin/windows-amd64.zip _bin/windows-amd64/dd-sesame.exe
	zip -9 _bin/linux-amd64.zip _bin/linux-amd64/dd-sesame
	zip -9 _bin/darwin-amd64.zip _bin/darwin-amd64/dd-sesame

test: fmt
	go test -v github.com/tintoy/dd-sesame

version:
	echo "package main\n\n// ProgramVersion is the current version of the DD Sesame tool.\nconst ProgramVersion = \"v0.1 (`git rev-parse HEAD`)\"" > ./version-info.go
