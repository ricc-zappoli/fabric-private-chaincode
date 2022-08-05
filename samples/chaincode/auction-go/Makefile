# Copyright 2019 Intel Corporation
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

TOP = ../../..
include $(TOP)/ecc_go/Makefile

CC_NAME ?= auction
BUILD_DIR := _build
LIB_DIR := _build/lib

build:
	@if [ ! -d $(LIB_DIR) ]; then \
		mkdir -p $(LIB_DIR); \
	fi

clean:
	rm -rf $(BUILD_DIR)