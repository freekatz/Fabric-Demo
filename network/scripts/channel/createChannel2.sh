#!/bin/bash

export PATH=${PWD}/../config:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="channel"}

. scripts/utils.sh

createChannelTx() {
	configtxgen -profile Org2Channel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}2.tx -channelID ${CHANNEL_NAME}2

    configtxgen -profile Org2Channel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPChannel2Anchor.tx -channelID ${CHANNEL_NAME}2 -asOrg Org2
}

createChannel() {
    docker exec cli20 bash ./scripts/channel/createChannelBlock.sh ${CHANNEL_NAME} 2
    
    docker cp cli20:/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}2.block ./channel-artifacts/

}

# joinChannel ORG
joinChannel() {
    docker exec cli20 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 2
}

setAnchorPeer() {
    docker exec cli20 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 2 2
}

## Create channeltx
infoln "Generating channel create transaction '${CHANNEL_NAME}2.tx'"
createChannelTx

## Create channel
infoln "Creating channel ${CHANNEL_NAME}2"
createChannel
successln "Channel '${CHANNEL_NAME}2' created"
sleep 3

## Join all the peers to the channel
infoln "Joining org2 peers to the ${CHANNEL_NAME}2..."
joinChannel
sleep 3

## Set the anchor peers for each org in the channel
infoln "Setting anchor peer for org2 in ${CHANNEL_NAME}2..."
setAnchorPeer
successln "Channel '${CHANNEL_NAME}2' joined"
