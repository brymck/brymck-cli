PROTOS := brymck/dates/v1/date brymck/calendar/v1/calendar_api brymck/risk/v1/risk_api brymck/securities/v1/securities_api

GO_FILES := $(shell find . -name '*.go')
PROTO_FILES := $(shell find proto -name '*.proto' 2>/dev/null) $(foreach proto,$(PROTOS),proto/$(proto).proto)
PROTO_PATH := /usr/local/include
GENPROTO_FILES := $(patsubst proto/%.proto,genproto/%.pb.go,$(PROTO_FILES))

all: proto test build

init: .init.stamp

.init.stamp:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go mod download
	touch $@

proto: $(PROTO_FILES) $(GENPROTO_FILES)

proto/brymck/dates/v1/date.proto:
	mkdir -p $(dir $@)
	curl --fail --location --output $@ --silent --show-error https://raw.githubusercontent.com/brymck/protos/master/brymck/dates/v1/date.proto
	echo >> $@
	echo 'option go_package = "github.com/brymck/brymck-cli/genproto/brymck/dates/v1";' >> $@

proto/brymck/calendar/v1/calendar_api.proto:
	mkdir -p $(dir $@)
	curl --fail --location --output $@ --silent --show-error https://$(GITHUB_TOKEN)@raw.githubusercontent.com/brymck/calendar-service/master/$@

proto/brymck/risk/v1/risk_api.proto:
	mkdir -p $(dir $@)
	curl --fail --location --output $@ --silent --show-error https://$(GITHUB_TOKEN)@raw.githubusercontent.com/brymck/risk-service/master/$@

proto/brymck/securities/v1/securities_api.proto:
	mkdir -p $(dir $@)
	curl --fail --location --output $@ --silent --show-error https://$(GITHUB_TOKEN)@raw.githubusercontent.com/brymck/securities-service/master/$@

genproto/%.pb.go: proto/%.proto | .init.stamp
	mkdir -p $(dir $@)
	protoc -Iproto -I$(PROTO_PATH) --go_out=paths=source_relative,plugins=grpc:genproto $<

test: profile.out

profile.out: $(GO_FILES) $(GENPROTO_FILES) | .init.stamp
	go test -race -coverprofile=profile.out -covermode=atomic ./...

build: blm

blm: $(GO_FILES) $(GENPROTO_FILES) | .init.stamp
	go build -ldflags='-w -s' -o $@ cmd/cli/*.go

clean:
	rm -rf proto/ genproto/ .init.stamp profile.out client service

.PHONY: all init proto test build clean
