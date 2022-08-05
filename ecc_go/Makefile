# Copyright 2019 Intel Corporation
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

include $(TOP)/build.mk

GOTAGS += -tags $(CC_NAME)
CAAS_PORT ?= 9999

# the following are the required docker build parameters
DOCKER_IMAGE ?= fpc/fpc-$(CC_NAME)-go${HW_EXTENSION}
HW_EXTENSION=$(shell if [ "${SGX_MODE}" = "HW" ]; then echo "-hw"; fi)

build: ecc docker

ecc: ecc_dependencies
	ego-go build $(GOTAGS) -o $(BUILD_DIR)/$(CC_NAME)-go $(TOP)/ecc_go/main.go
	cd $(BUILD_DIR) && ego sign
	cd $(BUILD_DIR) && ego uniqueid $(CC_NAME)-go > lib/mrenclave
	cd $(BUILD_DIR) && ego bundle $(CC_NAME)-go

ecc_dependencies:
	# hard to list explicitly, so just leave empty target,
	# which forces ecc to always be built


# Note:
# - docker images are not necessarily rebuild if they exist but are outdated.
#   To force rebuild you have two options
#   - do a 'make clobber' first. This ensures you will have the uptodate images
#     but is a broad and slow brush
#   - to just fore rebuilding an image, call `make` with DOCKER_FORCE_REBUILD defined
#   - to keep docker build quiet unless there is an error, call `make` with DOCKER_QUIET_BUILD defined
DOCKER_BUILD_OPTS ?=
ifdef DOCKER_QUIET_BUILD
	DOCKER_BUILD_OPTS += --quiet
endif
ifdef DOCKER_FORCE_REBUILD
	DOCKER_BUILD_OPTS += --no-cache
endif
ifeq ($(SGX_MODE), SIM)
    DOCKER_BUILD_OPTS += --build-arg OE_SIMULATION=1
endif
DOCKER_BUILD_OPTS += --build-arg FPC_VERSION=$(FPC_VERSION)
DOCKER_BUILD_OPTS += --build-arg SGX_MODE=$(SGX_MODE)
DOCKER_BUILD_OPTS += --build-arg CAAS_PORT=$(CAAS_PORT)
DOCKER_BUILD_OPTS += --build-arg CC_NAME=$(CC_NAME)
DOCKER_BUILD_OPTS += --build-arg BUILD_DIR=$(BUILD_DIR)



docker:
	$(DOCKER) build $(DOCKER_BUILD_OPTS) -t $(DOCKER_IMAGE):$(FPC_VERSION) -f $(TOP)/ecc_go/Dockerfile .
