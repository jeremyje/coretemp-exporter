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

# https://github.com/protocolbuffers/protobuf/releases
PROTOC_VERSION = 21.9

EXE = 
FX_FIND = find

ifeq ($(OS),Windows_NT)
	HOST_OS = windows
	HOST_PLATFORM = windows_amd64
	HOST_ARCH = amd64
	# Give priority to /usr/bin because it conflicts with C:\Windows\system32 within Msys32 environment.
	FX_FIND = /usr/bin/find.exe
	EXE = .exe
	SED_REPLACE = sed -i
	PROTOC_PACKAGE = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-win64.zip
else
	UNAME_S := $(shell uname -s)
	UNAME_ARCH := $(shell uname -m)
	ifeq ($(UNAME_S),Linux)
		HOST_OS = linux
		SED_REPLACE = sed -i
		ifeq ($(UNAME_ARCH),arm)
			HOST_PLATFORM = linux_arm
			HOST_ARCH = arm
			PROTOC_PACKAGE = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-aarch_64.zip
		else
			HOST_PLATFORM = linux_amd64
			HOST_ARCH = amd64
			PROTOC_PACKAGE = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
		endif
	endif
	ifeq ($(UNAME_S),Darwin)
		HOST_OS = darwin
		HOST_PLATFORM = darwin_amd64
		SED_REPLACE = sed -i ''
		PROTOC_PACKAGE = https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-osx-x86_64.zip
	endif
endif

REPOSITORY_ROOT := $(patsubst %/,%,$(dir $(abspath fx.mk)))
BUILD_DIR = $(REPOSITORY_ROOT)/build
ARCHIVES_DIR = $(BUILD_DIR)/archives
TOOLCHAIN_DIR = $(BUILD_DIR)/toolchain
TOOLCHAIN_BIN = $(TOOLCHAIN_DIR)/bin
THIRDPARTY_DIR = $(REPOSITORY_ROOT)/third_party

FX_GO = GO111MODULE=on go
FX_GO_INSTALL = GOPATH=$(TOOLCHAIN_DIR) $(FX_GO) install
FX_CURL = curl --retry 5 --retry-connrefused
PROTOC := GOPATH=$(TOOLCHAIN_DIR) $(TOOLCHAIN_BIN)/protoc
PROTOC_INCLUDE_FLAGS = -I $(REPOSITORY_ROOT) -I $(THIRDPARTY_DIR)/grpc_gateway/include/ -I $(THIRDPARTY_DIR)/google_protobuf/include/

PROTOC_TOOLCHAIN = build/toolchain/bin/protoc$(EXE)
PROTOC_TOOLCHAIN += build/toolchain/bin/protoc-gen-go$(EXE)
PROTOC_TOOLCHAIN += build/toolchain/bin/protoc-gen-go-grpc$(EXE)
PROTOC_TOOLCHAIN += build/toolchain/bin/protoc-gen-grpc-gateway$(EXE)
PROTOC_TOOLCHAIN += build/toolchain/bin/protoc-gen-openapiv2$(EXE)

build/toolchain/bin: $(PROTOC_TOOLCHAIN)
	touch $(TOOLCHAIN_BIN)

build/archives/protoc.zip:
	mkdir -p $(ARCHIVES_DIR)/
	$(FX_CURL) -o $(ARCHIVES_DIR)/protoc.zip -L $(PROTOC_PACKAGE)
	touch $@

build/toolchain/bin/protoc$(EXE): build/archives/protoc.zip
	mkdir -p $(TOOLCHAIN_BIN)/
	mkdir -p $(TOOLCHAIN_DIR)/protoc-temp/
	cp $(ARCHIVES_DIR)/protoc.zip $(TOOLCHAIN_DIR)/protoc-temp/
	(cd $(TOOLCHAIN_DIR)/protoc-temp/; unzip -q -o protoc.zip)
	cp $(TOOLCHAIN_DIR)/protoc-temp/bin/protoc$(EXE) $(TOOLCHAIN_BIN)/protoc$(EXE)
	rm -rf $(TOOLCHAIN_DIR)/protoc-temp/
	touch $@

build/toolchain/bin/protoc-gen-go$(EXE): third_party/google_protobuf/include/google/ third_party/grpc_gateway/include/protoc-gen-openapiv2/
	mkdir -p $(TOOLCHAIN_BIN)/
	cd $(TOOLCHAIN_BIN) && $(FX_GO_INSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	touch $@

build/toolchain/bin/protoc-gen-go-grpc$(EXE): third_party/google_protobuf/include/google/ third_party/grpc_gateway/include/protoc-gen-openapiv2/
	mkdir -p $(TOOLCHAIN_BIN)/
	cd $(TOOLCHAIN_BIN) && $(FX_GO_INSTALL) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	touch $@

build/toolchain/bin/protoc-gen-grpc-gateway$(EXE): third_party/google_protobuf/include/google/ third_party/grpc_gateway/include/protoc-gen-openapiv2/
	cd $(TOOLCHAIN_BIN) && $(FX_GO_INSTALL) github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	touch $@

build/toolchain/bin/protoc-gen-openapiv2$(EXE): third_party/google_protobuf/include/google/ third_party/grpc_gateway/include/protoc-gen-openapiv2/
	mkdir -p $(TOOLCHAIN_BIN)/
	cd $(TOOLCHAIN_BIN) && $(FX_GO_INSTALL) github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	touch $@

build/archives/googleapis.zip:
	mkdir -p $(ARCHIVES_DIR)/
	$(FX_CURL) -o $(ARCHIVES_DIR)/googleapis.zip -L \
		https://github.com/googleapis/googleapis/archive/master.zip
	touch $@

build/archives/grpc-gateway.zip:
	mkdir -p $(ARCHIVES_DIR)/
	$(FX_CURL) -o $(ARCHIVES_DIR)/grpc-gateway.zip -L \
		https://github.com/grpc-ecosystem/grpc-gateway/archive/master.zip
	touch $@

third_party/google_protobuf/include/google/: third_party/google_protobuf/include/google/LICENSE
	touch $@

third_party/google_protobuf/include/google/LICENSE: build/archives/protoc.zip build/archives/googleapis.zip
	rm -rf $(THIRDPARTY_DIR)/google_protobuf/include/
	mkdir -p $(TOOLCHAIN_DIR)/googleapis-temp/
	mkdir -p $(TOOLCHAIN_BIN)/
	mkdir -p $(THIRDPARTY_DIR)/google_protobuf/include/google/
	# Copy protobuf
	cp $(ARCHIVES_DIR)/protoc.zip $(TOOLCHAIN_DIR)/googleapis-temp/
	cp $(ARCHIVES_DIR)/googleapis.zip $(TOOLCHAIN_DIR)/googleapis-temp/
	(cd $(TOOLCHAIN_DIR)/googleapis-temp/; unzip -q -o protoc.zip)
	cp -rf $(TOOLCHAIN_DIR)/googleapis-temp/include/google/* \
		$(THIRDPARTY_DIR)/google_protobuf/include/google/

	# Copy google/apis
	(cd $(TOOLCHAIN_DIR)/googleapis-temp/; unzip -q -o googleapis.zip)
	mkdir -p $(THIRDPARTY_DIR)/google_protobuf/include/google/api/
	mkdir -p $(THIRDPARTY_DIR)/google_protobuf/include/google/rpc/
	mkdir -p $(THIRDPARTY_DIR)/google_protobuf/include/google/longrunning/
	cp -rf $(TOOLCHAIN_DIR)/googleapis-temp/googleapis-master/google/api/* \
		$(THIRDPARTY_DIR)/google_protobuf/include/google/api/
	cp -rf $(TOOLCHAIN_DIR)/googleapis-temp/googleapis-master/google/rpc/* \
		$(THIRDPARTY_DIR)/google_protobuf/include/google/rpc/
	cp -rf $(TOOLCHAIN_DIR)/googleapis-temp/googleapis-master/google/longrunning/* \
		$(THIRDPARTY_DIR)/google_protobuf/include/google/longrunning/
	cp -f $(TOOLCHAIN_DIR)/googleapis-temp/googleapis-master/LICENSE \
		$(THIRDPARTY_DIR)/google_protobuf/include/google/LICENSE
	$(FX_FIND) $(THIRDPARTY_DIR)/google_protobuf/include/google/ -type f -name '*BUILD.bazel' -exec rm {} +
	rm -rf $(TOOLCHAIN_DIR)/googleapis-temp
	touch $@

third_party/grpc_gateway/include/protoc-gen-openapiv2/: third_party/grpc_gateway/include/protoc-gen-openapiv2/LICENSE.txt
	touch $@

third_party/grpc_gateway/include/protoc-gen-openapiv2/LICENSE.txt: build/archives/grpc-gateway.zip
	rm -rf $(THIRDPARTY_DIR)/grpc_gateway/
	mkdir -p $(TOOLCHAIN_DIR)/grpc-gateway-temp/
	mkdir -p $(TOOLCHAIN_BIN)/
	mkdir -p $(THIRDPARTY_DIR)/grpc_gateway/include/protoc-gen-openapiv2/options/

	cp $(ARCHIVES_DIR)/grpc-gateway.zip $(TOOLCHAIN_DIR)/grpc-gateway-temp/
	(cd $(TOOLCHAIN_DIR)/grpc-gateway-temp/; unzip -q -o grpc-gateway.zip)
	cp -rf $(TOOLCHAIN_DIR)/grpc-gateway-temp/grpc-gateway-master/protoc-gen-openapiv2/options/*.proto \
		$(THIRDPARTY_DIR)/grpc_gateway/include/protoc-gen-openapiv2/options/
	cp -f $(TOOLCHAIN_DIR)/grpc-gateway-temp/grpc-gateway-master/LICENSE.txt \
		$(THIRDPARTY_DIR)/grpc_gateway/include/protoc-gen-openapiv2/LICENSE.txt
	$(FX_FIND) $(THIRDPARTY_DIR)/grpc_gateway/include/protoc-gen-openapiv2/ -type f -name '*BUILD.bazel' -exec rm {} +
	rm -rf $(TOOLCHAIN_DIR)/grpc-gateway-temp
	touch $@

no-sudo:
ifndef ALLOW_BUILD_WITH_SUDO
ifeq ($(shell whoami),root)
	@echo "ERROR: Running Makefile as root (or sudo)"
	@echo "Please follow the instructions at https://docs.docker.com/install/linux/linux-postinstall/ if you are trying to sudo run the Makefile because of the 'Cannot connect to the Docker daemon' error."
	@echo "NOTE: sudo/root do not have the authentication token to talk to any GCP service via gcloud."
	exit 1
endif
endif

%_grpc.pb.go: %.proto %.pb.go $(PROTOC_TOOLCHAIN)
	$(PROTOC) $(PROTOC_INCLUDE_FLAGS) --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
	$(FX_GO) fmt $@
	touch $@

%.pb.go: %.proto $(PROTOC_TOOLCHAIN)
	$(PROTOC) $(PROTOC_INCLUDE_FLAGS) --go_out=. --go_opt=paths=source_relative $<
	$(FX_GO) fmt $@
	touch $@

%.pb.gw.go: %.proto %_grpc.pb.go $(PROTOC_TOOLCHAIN)
	echo $(dir $<)
	$(PROTOC) $(PROTOC_INCLUDE_FLAGS) --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true --grpc-gateway_opt allow_delete_body=true $<
	$(SED_REPLACE) 's/proto_0/proto/g' $@
	$(SED_REPLACE) 's/status_0/status/g' $@
	$(FX_GO) fmt $@
	touch $@

%.swagger.json: %.proto %.pb.gw.go $(PROTOC_TOOLCHAIN)
	$(PROTOC) $(PROTOC_INCLUDE_FLAGS) --openapiv2_out . --openapiv2_opt logtostderr=true $<
	touch $@