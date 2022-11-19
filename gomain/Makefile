# Copyright 2022 Jeremy Edwards
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO = GO111MODULE=on go
DOCKER = DOCKER_CLI_EXPERIMENTAL=enabled docker

BASE_VERSION = 0.0.0-dev
SHORT_SHA = $(shell git rev-parse --short=7 HEAD | tr -d [:punct:])
VERSION_SUFFIX = $(SHORT_SHA)
VERSION = $(BASE_VERSION)-$(VERSION_SUFFIX)
BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
TAG := v$(VERSION)

export PATH := $(PWD)/build/toolchain/bin:$(PATH):/root/go/bin:/usr/local/go/bin:/usr/go/bin
GO = go
SOURCE_DIRS=$(shell go list ./... | grep -v '/vendor/')

PROTOS = proto/hardware.pb.go

ASSETS = $(PROTOS)
NICHE_PLATFORMS =
LINUX_PLATFORMS = linux_386 linux_amd64 linux_arm_v5 linux_arm_v6 linux_arm_v7 linux_arm64 linux_s390x linux_ppc64le linux_riscv64 linux_mips64le linux_mips linux_mipsle linux_mips64
LINUX_NICHE_PLATFORMS = 
WINDOWS_PLATFORMS = windows_386 windows_amd64
MAIN_PLATFORMS = windows_amd64 linux_amd64 linux_arm64
ALL_PLATFORMS = $(LINUX_PLATFORMS) $(LINUX_NICHE_PLATFORMS) $(WINDOWS_PLATFORMS) $(foreach niche,$(NICHE_PLATFORMS),$(niche)_amd64 $(niche)_arm64)
ALL_APPS = example

MAIN_BINARIES = $(foreach app,$(ALL_APPS),$(foreach platform,$(MAIN_PLATFORMS),build/bin/$(platform)/$(app)$(if $(findstring windows_,$(platform)),.exe,)))
ALL_BINARIES = $(foreach app,$(ALL_APPS),$(foreach platform,$(ALL_PLATFORMS),build/bin/$(platform)/$(app)$(if $(findstring windows_,$(platform)),.exe,)))

binaries: $(MAIN_BINARIES)
all: $(ALL_BINARIES)

build/bin/%:
	GOOS=$(firstword $(subst _, ,$(notdir $(abspath $(dir $@))))) GOARCH=$(word 2, $(subst _, ,$(notdir $(abspath $(dir $@))))) GOARM=$(subst v,,$(word 3, $(subst _, ,$(notdir $(abspath $(dir $@)))))) CGO_ENABLED=0 $(GO) build -o $@ cmd/$(basename $(notdir $@))/$(basename $(notdir $@)).go
	touch $@

lint:
	$(GO) fmt ./...
	$(GO) vet ./...

test:
	$(GO) test -race ${SOURCE_DIRS} -cover

coverage.txt:
	for sfile in ${SOURCE_DIRS} ; do \
		go test -race "$$sfile" -coverprofile=package.coverage -covermode=atomic; \
		if [ -f package.coverage ]; then \
			cat package.coverage >> coverage.txt; \
			$(RM) package.coverage; \
		fi; \
	done

clean:
	-chmod -R +w build/
	rm -rf build/

.PHONY: all run lint test clean
