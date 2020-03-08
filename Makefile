all: test build

test: profile.out

profile.out: $(GO_FILES)
	go mod download
	go test -race -coverprofile=profile.out -covermode=atomic ./...

build: blm

blm: $(GO_FILES)
	go mod download
	go build -ldflags='-w -s' -o $@ cmd/cli/*.go

clean:
	rm -rf profile.out blm

.PHONY: all test build
