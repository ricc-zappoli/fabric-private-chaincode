/*
Copyright IBM Corp. All Rights Reserved.
Copyright 2020 Intel Corporation

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/hyperledger/fabric-private-chaincode/samples/application/simple-cli-go/cmd"
)

const FileName = "goSIM.csv"

func main() {
	start := time.Now()
	cmd.Execute()
	end := time.Since(start)

	csvFile, _ := os.OpenFile(FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w := csv.NewWriter(csvFile)
	w.Write(append([]string{
		time.Now().UTC().Format("2006-01-02 15:04:05.000"),
		end.String(),
	}, os.Args[1:]...))
	w.Flush()

}
