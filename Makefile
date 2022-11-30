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

include proto.mk

GO = GO111MODULE=on go
DOCKER = DOCKER_CLI_EXPERIMENTAL=enabled docker

VERSION = $(shell git describe --tags)
BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
TAG := $(VERSION)

export PATH := $(PWD)/build/toolchain/bin:$(PATH):/root/go/bin:/usr/local/go/bin:/usr/go/bin
GO = go
SOURCE_DIRS=$(shell go list ./... | grep -v '/vendor/')

REGISTRY = ghcr.io/jeremyje
CORETEMP_EXPORTER_IMAGE = $(REGISTRY)/coretemp-exporter

PROTOS = proto/hardware.pb.go

ASSETS = $(PROTOS)
NICHE_PLATFORMS =
LINUX_PLATFORMS = linux_386 linux_amd64 linux_arm_v5 linux_arm_v6 linux_arm_v7 linux_arm64 linux_s390x linux_ppc64le linux_riscv64 linux_mips64le linux_mips linux_mipsle linux_mips64
LINUX_NICHE_PLATFORMS = 
WINDOWS_PLATFORMS = windows_386 windows_amd64
MAIN_PLATFORMS = windows_amd64 linux_amd64 linux_arm64
ALL_PLATFORMS = $(LINUX_PLATFORMS) $(LINUX_NICHE_PLATFORMS) $(WINDOWS_PLATFORMS) $(foreach niche,$(NICHE_PLATFORMS),$(niche)_amd64 $(niche)_arm64)
ALL_APPS = coretemp-exporter

MAIN_BINARIES = $(foreach app,$(ALL_APPS),$(foreach platform,$(MAIN_PLATFORMS),build/bin/$(platform)/$(app)$(if $(findstring windows_,$(platform)),.exe,)))
ALL_BINARIES = $(foreach app,$(ALL_APPS),$(foreach platform,$(ALL_PLATFORMS),build/bin/$(platform)/$(app)$(if $(findstring windows_,$(platform)),.exe,)))

WINDOWS_VERSIONS = 1709 1803 1809 1903 1909 2004 20H2 ltsc2022
BUILDX_BUILDER = buildx-builder

binaries: $(MAIN_BINARIES)
all: $(ALL_BINARIES)
assets: $(ASSETS)
protos: $(PROTOS)

build/bin/%: $(ASSETS)
	GOOS=$(firstword $(subst _, ,$(notdir $(abspath $(dir $@))))) GOARCH=$(word 2, $(subst _, ,$(notdir $(abspath $(dir $@))))) GOARM=$(subst v,,$(word 3, $(subst _, ,$(notdir $(abspath $(dir $@)))))) CGO_ENABLED=0 $(GO) build -o $@ cmd/$(basename $(notdir $@))/$(basename $(notdir $@)).go
	touch $@

run: cmd/coretemp-exporter/coretemp-exporter.go
	$(GO) run cmd/coretemp-exporter/coretemp-exporter.go -log=cputemps.ndjson -endpoint=:8181

run8081: cmd/coretemp-exporter/coretemp-exporter.go
	$(GO) run cmd/coretemp-exporter/coretemp-exporter.go -log=cputemps.ndjson

lint:
	$(GO) fmt ./...
	$(GO) vet ./...

test:
	$(GO) test -race ${SOURCE_DIRS} -cover -count 50

test-25:
	$(GO) test -race ${SOURCE_DIRS} -cover -count 25

coverage.txt:
	for sfile in ${SOURCE_DIRS} ; do \
		go test -race "$$sfile" -coverprofile=package.coverage -covermode=atomic; \
		if [ -f package.coverage ]; then \
			cat package.coverage >> coverage.txt; \
			$(RM) package.coverage; \
		fi; \
	done

ensure-builder:
	-$(DOCKER) buildx create --name $(BUILDX_BUILDER)

ALL_IMAGES = $(CORETEMP_EXPORTER_IMAGE) $(CERTTOOL_IMAGE) $(HTTPPROBE_IMAGE)
# https://github.com/docker-library/official-images#architectures-other-than-amd64
images: DOCKER_PUSH = --push
images: linux-images windows-images
	-$(DOCKER) manifest rm $(CORETEMP_EXPORTER_IMAGE):$(TAG)

	for image in $(ALL_IMAGES) ; do \
		$(DOCKER) manifest create $$image:$(TAG) $(foreach winver,$(WINDOWS_VERSIONS),$${image}:$(TAG)-windows_amd64-$(winver)) $(foreach platform,$(LINUX_PLATFORMS),$${image}:$(TAG)-$(platform)) ; \
		for winver in $(WINDOWS_VERSIONS) ; do \
			windows_version=`$(DOCKER) manifest inspect mcr.microsoft.com/windows/nanoserver:$${winver} | jq -r '.manifests[0].platform["os.version"]'`; \
			$(DOCKER) manifest annotate --os-version $${windows_version} $${image}:$(TAG) $${image}:$(TAG)-windows_amd64-$${winver} ; \
		done ; \
		$(DOCKER) manifest push $$image:$(TAG) ; \
	done

ALL_LINUX_IMAGES = $(foreach app,$(ALL_APPS),$(foreach platform,$(LINUX_PLATFORMS),linux-image-$(app)-$(platform)))
linux-images: $(ALL_LINUX_IMAGES)

linux-image-coretemp-exporter-%: build/bin/%/coretemp-exporter ensure-builder
	$(DOCKER) buildx build --builder $(BUILDX_BUILDER) --platform $(subst _,/,$*) --build-arg BINARY_PATH=$< -f cmd/coretemp-exporter/Dockerfile -t $(CORETEMP_EXPORTER_IMAGE):$(TAG)-$* . $(DOCKER_PUSH)

ALL_WINDOWS_IMAGES = $(foreach app,$(ALL_APPS),$(foreach winver,$(WINDOWS_VERSIONS),windows-image-$(app)-$(winver)))
windows-images: $(ALL_WINDOWS_IMAGES)

windows-image-coretemp-exporter-%: build/bin/windows_amd64/coretemp-exporter.exe ensure-builder
	$(DOCKER) buildx build --builder $(BUILDX_BUILDER) --platform windows/amd64 -f cmd/coretemp-exporter/Dockerfile.windows --build-arg WINDOWS_VERSION=$* -t $(CORETEMP_EXPORTER_IMAGE):$(TAG)-windows_amd64-$* . $(DOCKER_PUSH)

clean:
	rm -f coverage.txt
	-chmod -R +w build/
	rm -rf build/

convert:
	$(GO) run cmd/converter/converter.go

.PHONY: all run lint test clean
