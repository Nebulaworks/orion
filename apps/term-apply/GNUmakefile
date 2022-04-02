COMMIT = $$(git rev-parse HEAD)
BUILD_TIME = $$(date +'%Y-%m-%d_%T')
VERSION_PATH = "github.com/nebulaworks/orion/apps/term-apply/pkg/version"

build:	*.go
	go build -ldflags="-X '${VERSION_PATH}.Commit=${COMMIT}' -X '${VERSION_PATH}.BuildTime=${BUILD_TIME}'" -o term-apply

term-apply: build

deps:
	go mod verify
	go mod tidy

test:
	go clean -testcache
	go test --race ./...
