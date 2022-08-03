#!/bin/bash

cd $FPC_PATH/samples/application/simple-cli-go
make

for i in {1..500}; do

    cd $FPC_PATH/samples/deployment/test-network
    /bin/bash restart.sh -d

    cd $FPC_PATH/samples/application/simple-cli-go
    export CC_NAME=${CC_ID}
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
    sleep 10

    if [ "${SGX_MODE}" = "HW" ]; then 
        hw_mode = true
        export SGX_MODE=SIM
    fi

    ./fpcclient init $CORE_PEER_ID
    ./fpcclient invoke init House$i
    ./fpcclient invoke create Auction
    ./fpcclient invoke submit Auction John 100
    ./fpcclient invoke submit Auction Jane 200
    ./fpcclient invoke submit Auction John 400
    ./fpcclient invoke submit Auction Danny 300
    ./fpcclient invoke close Auction
    ./fpcclient invoke eval Auction

    if [ "$hw_mode" = true ]; then 
        export SGX_MODE=HW
    fi

done

cd $FPC_PATH/samples/deployment/test-network
/bin/bash stop.sh
