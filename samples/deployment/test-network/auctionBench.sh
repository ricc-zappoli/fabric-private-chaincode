#!/bin/bash

/bin/bash restart.sh -d

# Chaincode execution
#------------
cd $FPC_PATH/samples/deployment/test-network
./update-connection.sh

cd $FPC_PATH/samples/application/simple-cli-go
make

export CC_NAME=${CC_ID}
echo CC_NAME=${CC_NAME}
export CHANNEL_NAME=mychannel
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_ID=peer0.org1.example.com
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=$FPC_PATH/samples/deployment/test-network/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_CERT_FILE=$FPC_PATH/samples/deployment/test-network/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
export CORE_PEER_TLS_ENABLED="true"
export CORE_PEER_TLS_KEY_FILE=$FPC_PATH/samples/deployment/test-network/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE=$FPC_PATH/samples/deployment/test-network/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORDERER_CA=$FPC_PATH/samples/deployment/test-network/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export GATEWAY_CONFIG=$FPC_PATH/samples/deployment/test-network/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.yaml

echo "Waiting for chaincodes..."
sleep 15

# TODO: switch to SIM mode if HW to avoid attestation checks

./fpcclient init $CORE_PEER_ID
./fpcclient invoke init House
./fpcclient invoke create Auction
./fpcclient invoke submit Auction John 3
./fpcclient invoke close Auction
./fpcclient invoke eval Auction

# TODO: switch back to HW

