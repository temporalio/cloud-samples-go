$(VERBOSE).SILENT:
############################# Main targets #############################
# Install dependencies.
install: buf-install

# Run all linters and compile proto files.
proto: copy-api-cloud-api grpc
########################################################################

##### Variables ######
ifndef GOPATH
GOPATH := $(shell go env GOPATH)
endif

GOBIN := $(if $(shell go env GOBIN),$(shell go env GOBIN),$(GOPATH)/bin)
SHELL := PATH=$(GOBIN):$(PATH) /bin/sh

COLOR := "\e[1;36m%s\e[0m\n"

PROTO_ROOT := proto
PROTO_OUT := protogen
$(PROTO_OUT):
	mkdir $(PROTO_OUT)

##### Copy the proto files from the api-cloud repo #####
copy-api-cloud-api:
	@printf $(COLOR) "Copy api-cloud..."
	rm -rf $(PROTO_ROOT)
	mkdir -p $(PROTO_ROOT)/temporal/api
	git clone git@github.com:temporalio/api-cloud.git --depth=1 --branch abhinav/userMgmt --single-branch $(PROTO_ROOT)/api-cloud-tmp
	mv -f $(PROTO_ROOT)/api-cloud-tmp/temporal/api/cloud $(PROTO_ROOT)/temporal/api
	rm -rf $(PROTO_ROOT)/api-cloud-tmp

##### Compile proto files for go #####
grpc: go-grpc fix-proto-generated-go-path

go-grpc: clean $(PROTO_OUT)
	printf $(COLOR) "Compile for go-gRPC..."
	cd proto && buf generate --output ../

fix-proto-generated-go-path:
	@if [ "$$(uname -s)" = "Darwin" ]; then find $(PROTO_OUT) -name '*.go' -exec sed -i '' "s/go\.temporal\.io/github\.com\/temporalio\/cloud-operations-workflows\/protogen\/temporal/g" {} \;; else find $(PROTO_OUT) -name '*.go' -exec sed -i 's/go\.temporal\.io/github\.com\/temporalio\/cloud-operations-workflows\/protogen\/temporal/g' {} \;; fi

##### Plugins & tools #####
buf-install:
	printf $(COLOR) "Install/update buf..."
	go install github.com/bufbuild/buf/cmd/buf@v1.25.1

##### Clean #####
clean:
	printf $(COLOR) "Delete generated go files..."
	rm -rf $(PROTO_OUT)
