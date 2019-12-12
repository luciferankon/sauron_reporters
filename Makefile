PKGS := $(shell go  list ./... | grep -v /vendor)

.PHONY: test
test:
	go test -cover $(PKGS)

reporter:
	CGO_ENABLED=0 go build -o bin/reporter ./pkg/main/

.PHONY: reporter_stripped
reporter_stripped:
	go build -o bin/reporter -ldflags="-s -w" ./pkg/main/

.PHONY: reporter_compressed
reporter_compressed: reporter_stripped
	upx bin/reporter