#!/bin/bash

# Copyright IBM Corp. All Rights Reserved.
# Copyright 2020 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

make -C $FPC_PATH/samples/deployment/test-network ercc-ecc-stop

cd $FPC_PATH/samples/deployment/test-network/fabric-samples/test-network
./network.sh down

cd $FPC_PATH/samples/deployment/test-network
