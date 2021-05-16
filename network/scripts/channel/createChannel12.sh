#!/bin/bash

export PATH=${PWD}/../config:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="channel"}

. scripts/utils.sh

createChannelTx() {
	configtxgen -profile Org12Channel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}12.tx -channelID ${CHANNEL_NAME}12

    configtxgen -profile Org12Channel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPChannel12Anchor.tx -channelID ${CHANNEL_NAME}12 -asOrg Org1

    configtxgen -profile Org12Channel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPChannel12Anchor.tx -channelID ${CHANNEL_NAME}12 -asOrg Org2
}

createChannel() {
    docker exec cli10 bash ./scripts/channel/createChannelBlock.sh ${CHANNEL_NAME} 12
    
    docker cp cli10:/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}12.block ./channel-artifacts/

    docker cp ./channel-artifacts/${CHANNEL_NAME}12.block cli11:/opt/gopath/src/github.com/hyperledger/fabric/peer/ 

    docker cp ./channel-artifacts/${CHANNEL_NAME}12.block cli20:/opt/gopath/src/github.com/hyperledger/fabric/peer/ 
}

# joinChannel ORG
joinChannel() {
    docker exec cli10 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 12
    docker exec cli20 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 12
}

setAnchorPeer() {
    docker exec cli10 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 12 1
    docker exec cli20 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 12 2
}

## Create channeltx
infoln "Generating channel create transaction '${CHANNEL_NAME}12.tx'"
createChannelTx

## Create channel
infoln "Creating channel ${CHANNEL_NAME}12"
createChannel
successln "Channel '${CHANNEL_NAME}12' created"
sleep 3

## Join all the peers to the channel
infoln "Joining org1 and org2 peers to the ${CHANNEL_NAME}12..."
joinChannel
sleep 3

## Set the anchor peers for each org in the channel
infoln "Setting anchor peer for org1 and org2 in ${CHANNEL_NAME}12..."
setAnchorPeer
successln "Channel '${CHANNEL_NAME}12' joined"
