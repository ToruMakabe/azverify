APP=azverify

.PHONY: test
test:
	go test -v -count=1 -race ./...

.PHONY: clean
clean:
	go clean

.PHONY: build
build: clean
	go build -o ${APP} ./main.go

.PHONY: release-test
release-test:
	goreleaser --snapshot --skip-publish --rm-dist
