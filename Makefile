(VERBOSE).SILENT:
############################# Main targets #############################
# Build the binaries.
bins: clean worker exporttool

########################################################################

##### Variables ######
ifndef GOPATH
GOPATH := $(shell go env GOPATH)
endif

GOBIN := $(if $(shell go env GOBIN),$(shell go env GOBIN),$(GOPATH)/bin)
SHELL := PATH=$(GOBIN):$(PATH) /bin/sh

COLOR := "\e[1;36m%s\e[0m\n"

bins: clean worker exporttool

##### Build #####
.PHONY: worker
worker:
	@rm -rf ./worker
	@go build -o worker ./cmd/worker/*.go

.PHONY: exporttool
exporttool:
	@rm -rf ./exporttool
	@go build -o exporttool ./cmd/exporttool/*.go

clean:
	@rm -rf ./worker
	@rm -rf ./exporttool

