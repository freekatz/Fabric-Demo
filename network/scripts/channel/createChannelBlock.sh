#!/bin/bash

CHANNEL_NAME="$1"
CHANNEL_ID="$2"
: ${CHANNEL_NAME:="channel"}
: ${CHANNEL_ID:="1"} 

peer channel create -o orderer.example.com:7050 -c ${CHANNEL_NAME}${CHANNEL_ID} -f ./channel-artifacts/${CHANNEL_NAME}${CHANNEL_ID}.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/orgs/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
