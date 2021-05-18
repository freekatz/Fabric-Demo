#!/bin/bash

export PATH=${PWD}/../config:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="channel"}

. scripts/utils.sh

createChannelTx() {
	configtxgen -profile Org3Channel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}3.tx -channelID ${CHANNEL_NAME}3

    configtxgen -profile Org3Channel -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPChannel3Anchor.tx -channelID ${CHANNEL_NAME}3 -asOrg Org3
}

createChannel() {
    docker exec cli30 bash ./scripts/channel/createChannelBlock.sh ${CHANNEL_NAME} 3
    
    docker cp cli30:/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}3.block ./channel-artifacts/

}

# joinChannel ORG
joinChannel() {
    docker exec cli30 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 3
}

setAnchorPeer() {
    docker exec cli30 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 3 3
}

## Create channeltx
infoln "Generating channel create transaction '${CHANNEL_NAME}3.tx'"
createChannelTx

## Create channel
infoln "Creating channel ${CHANNEL_NAME}3"
createChannel
successln "Channel '${CHANNEL_NAME}3' created"
sleep 3

## Join all the peers to the channel
infoln "Joining org3 peers to the ${CHANNEL_NAME}3..."
joinChannel
sleep 3

## Set the anchor peers for each org in the channel
infoln "Setting anchor peer for org3 in ${CHANNEL_NAME}3..."
setAnchorPeer
successln "Channel '${CHANNEL_NAME}3' joined"
