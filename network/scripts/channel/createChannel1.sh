#!/bin/bash

export PATH=${PWD}/../config:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="channel"}

. scripts/utils.sh

createChannelTx() {
	configtxgen -profile Org1Channel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}1.tx -channelID ${CHANNEL_NAME}1

    configtxgen -profile Org1Channel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPChannel1Anchor.tx -channelID ${CHANNEL_NAME}1 -asOrg Org1
}

createChannel() {
    docker exec cli10 bash ./scripts/channel/createChannelBlock.sh ${CHANNEL_NAME} 1
    
    docker cp cli10:/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}1.block ./channel-artifacts/

}

# joinChannel ORG
joinChannel() {
    docker exec cli10 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 1
}

setAnchorPeer() {
    docker exec cli10 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 1 1
}

## Create channeltx
infoln "Generating channel create transaction '${CHANNEL_NAME}1.tx'"
createChannelTx

## Create channel
infoln "Creating channel ${CHANNEL_NAME}1"
createChannel
successln "Channel '${CHANNEL_NAME}1' created"
sleep 3

## Join all the peers to the channel
infoln "Joining org1 peers to the ${CHANNEL_NAME}1..."
joinChannel
sleep 3

## Set the anchor peers for each org in the channel
infoln "Setting anchor peer for org1 in ${CHANNEL_NAME}1..."
setAnchorPeer
successln "Channel '${CHANNEL_NAME}1' joined"
