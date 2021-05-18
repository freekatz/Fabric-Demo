#!/bin/bash

export PATH=${PWD}/../config:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="channel"}

. scripts/utils.sh

createChannelTx() {
	configtxgen -profile Org123Channel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}123.tx -channelID ${CHANNEL_NAME}123

    configtxgen -profile Org123Channel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPChannel123Anchor.tx -channelID ${CHANNEL_NAME}123 -asOrg Org1

    configtxgen -profile Org123Channel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPChannel123Anchor.tx -channelID ${CHANNEL_NAME}123 -asOrg Org2

    configtxgen -profile Org123Channel -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPChannel123Anchor.tx -channelID ${CHANNEL_NAME}123 -asOrg Org3
}

createChannel() {
    docker exec cli10 bash ./scripts/channel/createChannelBlock.sh ${CHANNEL_NAME} 123
    
    docker cp cli10:/opt/gopath/src/github.com/hyperledger/fabric/peer/${CHANNEL_NAME}123.block ./channel-artifacts/

    docker cp ./channel-artifacts/${CHANNEL_NAME}123.block cli20:/opt/gopath/src/github.com/hyperledger/fabric/peer/ 

    docker cp ./channel-artifacts/${CHANNEL_NAME}123.block cli30:/opt/gopath/src/github.com/hyperledger/fabric/peer/ 

}

# joinChannel ORG
joinChannel() {
    docker exec cli10 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 123
    docker exec cli20 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 123
    docker exec cli30 bash ./scripts/channel/joinChannel.sh ${CHANNEL_NAME} 123
}

setAnchorPeer() {
    docker exec cli10 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 123 1
    docker exec cli20 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 123 2
    docker exec cli30 bash ./scripts/channel/setAnchorPeer.sh ${CHANNEL_NAME} 123 3
}

## Create channeltx
infoln "Generating channel create transaction '${CHANNEL_NAME}123.tx'"
createChannelTx

## Create channel
infoln "Creating channel ${CHANNEL_NAME}123"
createChannel
successln "Channel '${CHANNEL_NAME}123' created"
sleep 3

## Join all the peers to the channel
infoln "Joining orgs peers to the ${CHANNEL_NAME}123..."
joinChannel
sleep 3

## Set the anchor peers for each org in the channel
infoln "Setting anchor peer for orgs in ${CHANNEL_NAME}123..."
setAnchorPeer
successln "Channel '${CHANNEL_NAME}123' joined"
