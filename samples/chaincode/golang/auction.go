/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package golang

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

const OK = "OK"
const AUCTION_DRAW = "DRAW"
const AUCTION_NO_BIDS = "NO_BIDS"
const AUCTION_ALREADY_EXISTING = "AUCTION_ALREADY_EXISTING"
const AUCTION_NOT_EXISTING = "AUCTION_NOT_EXISTING"
const AUCTION_ALREADY_CLOSED = "AUCTION_ALREADY_CLOSED"
const AUCTION_STILL_OPEN = "AUCTION_STILL_OPEN"

const INITIALIZED_KEY = "initialized"
const AUCTION_HOUSE_NAME_KEY = "auction_house_name"

const SEP = ".";
const PREFIX = SEP + "somePrefix" + SEP;

type Auction struct {
}

type auctionType struct {
	name string
	isOpen bool
}

type bidType struct {
	bidderName string
	value int
}

func (t *Auction) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *Auction) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	var bool initialized

	init, err = stub.GetState(INITIALIZED_KEY)
	if err != nil || !init {
		initialized := false
		auctionHouseName := "(uninitialized)"
	} else {
		ahn, err = stub.GetState(AUCTION_HOUSE_NAME_KEY)
		if err != nil {
			auctionHouseName := "(uninitialized)"
		} else {
			auctionHouseName := ahn
		}
	}

	fmt.Println("AuctionCC: +++ Executing", auctionHouseName, "auction chaincode invocation +++")
	functionName, params := stub.GetFunctionAndParameters()
	fmt.Println("AuctionCC: Function:", functionName, "Params:", params)

	if !initialized && functionName != "init" {
		return shim.Error("AuctionCC: Auction not yet initialized / No re-initialized allowed")
	}

	switch functionName {
	case "init":
		result := t.initAuctionHouse(stub, params[0])
	case "create":
		result := t.auctionCreate(stub, params[0])
	case "submit":
		result := t.auctionSubmit(stub, params[0], params[1], params[2])
	case "close":
		result := t.auctionClose(stub, params[0])
	case "eval":
		result := t.auctionEval(stub, params[0])
	default:
		result := shim.Error("AuctionCC: RECEIVED UNKOWN transaction")
	}

	fmt.Println("AuctionCC: Response:", result)
	fmt.Println("AuctionCC: +++ Executing done +++")
	return shim.Success([]byte(result))

}

func (t *Auction) initAuctionHouse(stub shim.ChaincodeStubInterface, auctionName) string {
	stub.PutState(AUCTION_HOUSE_NAME_KEY, auctionName)
	stub.PutState(INITIALIZED_KEY, 1)
	reutnr "OK"
}

func (t *Auction) auctionCreate(stub shim.ChaincodeStubInterface, auctionName) string {
	_, err := stub.GetState(auctionName)
	if err == nil {
		fmt.Println("AuctionCC: Auction already exists")
		return AUCTION_ALREADY_EXISTING
	}

	auction := &auctionType{name: auctionName, isOpen: true}

	auctionBytes, _ := json.Marshal(auction)
	stub.PutState(auctionName, auctionBytes)
	return "OK"
}

func (t *Auction) auctionSubmit(stub shim.ChaincodeStubInterface, auctionName, bidderName, value) string {
	auctionBytes, err := stub.GetState(auctionName)
	if err != nil {
		fmt.Println("AuctionCC: Auction does not exist")
		return AUCTION_NOT_EXISTING
	}

	var auction auctionType
	json.Unmarshal(auctionBytes, &auction)

	if !auction.isOpen {
		fmt.Println("AuctionCC: Auction is already closed")
		return AUCTION_ALREADY_CLOSED
	}

	key := PREFIX + auctionName + SEP + bidderName + SEP
	bid := &bidType{bidderName: bidderName, value: value}

	bidBytes, _ := json.Marshal(bid)
	stub.PutState(key, bid)
	return "OK"
}

func (t *Auction) auctionClose(stub shim.ChaincodeStubInterface, auctionName) string {
	auctionBytes, err := stub.GetState(auctionName)
	if err != nil {
		fmt.Println("AuctionCC: Auction does not exist")
		return AUCTION_NOT_EXISTING
	}

	var auction auctionType
	json.Unmarshal(auctionBytes, &auction)

	if !auction.isOpen {
		fmt.Println("AuctionCC: Auction is already closed")
		return AUCTION_ALREADY_CLOSED
	}

	auction.isOpen = false

	auctionBytes, _ := json.Marshal(auction)
	stub.PutState(auctionName, auctionBytes)
	return "OK"
}

func (t *Auction) auctionEval(stub shim.ChaincodeStubInterface, auctionName) string {
	auctionBytes, err := stub.GetState(auctionName)
	if err != nil {
		fmt.Println("AuctionCC: Auction does not exist")
		return AUCTION_NOT_EXISTING
	}

	var auction auctionType
	json.Unmarshal(auctionBytes, &auction)

	if auction.isOpen {
		fmt.Println("AuctionCC: Auction is still open")
		return AUCTION_STILL_OPEN
	}

	// TODO
}