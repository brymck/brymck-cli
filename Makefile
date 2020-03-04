PROTOS := securities_api

GO_FILES := $(shell find . -name '*.go')
PROTO_FILES := $(shell find proto -name '*.proto' 2>/dev/null) $(foreach proto,$(PROTOS),proto/$(proto).proto)
PROTO_PATH := /usr/local/include
GENPROTO_FILES := $(patsubst proto/%.proto,genproto/%.pb.go,$(PROTO_FILES))
GITHUB_TOKEN := $(GITHUB_TOKEN)

all: generate test build

init: .init.stamp

.init.stamp:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go mod download
	touch $@

generate: $(GENPROTO_FILES)

proto genproto:
	mkdir $@

proto/securities_api.proto: | proto
	curl --fail --location --output $@ --silent --show-error https://$(GITHUB_TOKEN)@raw.githubusercontent.com/brymck/securities-service/master/$@

genproto/%.pb.go: proto/%.proto | .init.stamp genproto
	protoc -Iproto -I$(PROTO_PATH) --go_out=plugins=grpc:$(dir $@) $<

test: profile.out

profile.out: $(GO_FILES) $(GENPROTO_FILES) | .init.stamp
	go test -race -coverprofile=profile.out -covermode=atomic ./...

build: blm

blm: $(GO_FILES) $(GENPROTO_FILES) | .init.stamp
	go build -ldflags='-w -s' -o $@ cmd/cli/*.go

clean:
	rm -rf proto/ genproto/ .init.stamp profile.out client service

.PHONY: all init generate test build clean
