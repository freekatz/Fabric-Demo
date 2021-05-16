#!/bin/bash

CHANNEL_NAME="$1"
CHANNEL_ID="$2"
ORG_ID="$3"
: ${CHANNEL_NAME:="channel"}
: ${CHANNEL_ID:="1"}
: ${ORG_ID:="1"}

peer channel update -o orderer.example.com:7050 -c ${CHANNEL_NAME}${CHANNEL_ID} -f ./channel-artifacts/Org${ORG_ID}MSPChannel${CHANNEL_ID}Anchor.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/orgs/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
