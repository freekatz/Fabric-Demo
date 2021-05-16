#!/bin/bash

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx
export VERBOSE=false

. scripts/utils.sh

function clearContainers() {
    CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')
    if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
        infoln "No containers available for deletion"
    else
        docker rm -f $CONTAINER_IDS
    fi
}

function removeUnwantedImages() {
    DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*/) {print $3}')
    if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
        infoln "No images available for deletion"
    else
        docker rmi -f $DOCKER_IMAGE_IDS
    fi
}

function genCrypto() {
    local ORG=$1
    infoln "Creating ${ORG} Identities"
    set -x
    cryptogen generate --config=./orgs/cryptogen/crypto-config-${ORG}.yaml --output="./orgs"
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        fatalln "Failed to generate ${ORG} certificates..."
    fi
}

function createOrgs() {
    if [ -d "orgs/peerOrganizations" ]; then
        rm -Rf orgs/peerOrganizations && rm -Rf orgs/ordererOrganizations
    fi

    # Create crypto material using cryptogen
    infoln "Generating certificates using cryptogen tool"

    genCrypto org1
    genCrypto org2
    genCrypto orderer

    infoln "Generating CCP files for Orgs"
    infoln "Please add extra peers config manually"
    bash orgs/ccp/ccp-generate.sh
}


function createConsortium() {
    infoln "Generating Orderer Genesis block"
    set -x
    configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID system-channel
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        fatalln "Failed to generate orderer genesis block..."
    fi
}

function networkUp() {
    # generate artifacts if they don't exist
    if [ ! -d "orgs/peerOrganizations" ]; then
        createOrgs
        createConsortium
    fi

    docker-compose up -d
}

function createChannel() {
    # Bring up the network if it is not already up.

    if [ ! -d "orgs/peerOrganizations" ]; then
        errorln "Have not up the network! Please first run the 'bash net.sh up'"
        exit
    fi
    
    bash scripts/channel/createChannel.sh $CHANNEL_NAME
}

## Call the script to deploy a chaincode to the channel
function deployCC() {
    bash scripts/cc/deployCC.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE $CC_MODE

    if [ $? -ne 0 ]; then
        fatalln "Deploying chaincode failed"
    fi
}

function networkStart() {
  docker-compose up -d
}

function networkStop() {
  docker-compose down -v
}

# Tear down running network
function networkDown() {
  set -x
  docker-compose down --volumes --remove-orphans
  # Bring down the network, deleting the volumes
  #Cleanup the chaincode containers
  clearContainers
  #Cleanup images
  removeUnwantedImages
  
  rm -rf orgs/peerOrganizations
  rm -rf orgs/ordererOrganizations
  rm -rf ../chaincode/releases/*.tar.gz
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

# channel name defaults to "channel"
CHANNEL_NAME="channel"
# channel id defaults to "channel"
CHANNEL_ID="1"
# chaincode name defaults to "NA"
CC_NAME="NA"
# chaincode path defaults to "NA"
CC_SRC_PATH="NA"
# endorsement policy defaults to "NA". This would allow chaincodes to use the majority default policy.
CC_END_POLICY="NA"
# collection configuration defaults to "NA"
CC_COLL_CONFIG="NA"
# chaincode init function defaults to "NA"
CC_INIT_FCN="NA"
# chaincode language defaults to "NA"
CC_SRC_LANGUAGE="NA"
# Chaincode version
CC_VERSION="1.0"
# Chaincode definition sequence
CC_SEQUENCE=1

CC_MODE=all

# Parse commandline args

## Parse mode
if [[ $# -lt 1 ]] ; then
  printHelp
  exit 0
else
  MODE=$1
  shift
fi

# parse a createChannel subcommand if used
if [[ $# -ge 1 ]] ; then
  key="$1"
  if [[ "$key" == "createChannel" ]]; then
      export MODE="createChannel"
      shift
  fi
fi

# parse flags

while [[ $# -ge 1 ]] ; do
  key="$1"
  case $key in
  -h )
    printHelp $MODE
    exit 0
    ;;
  -c )
    CHANNEL_NAME="$2"
    shift
    ;;
  -cid )
    CHANNEL_ID="$2"
    shift
    ;;
  -ccl )
    CC_SRC_LANGUAGE="$2"
    shift
    ;;
  -ccn )
    CC_NAME="$2"
    shift
    ;;
  -ccv )
    CC_VERSION="$2"
    shift
    ;;
  -ccs )
    CC_SEQUENCE="$2"
    shift
    ;;
  -ccp )
    CC_SRC_PATH="$2"
    shift
    ;;
  -ccep )
    CC_END_POLICY="$2"
    shift
    ;;
  -cccg )
    CC_COLL_CONFIG="$2"
    shift
    ;;
  -cci )
    CC_INIT_FCN="$2"
    shift
    ;;
  -ccm )
    CC_MODE="$2"
    shift
    ;;
  -verbose )
    VERBOSE=true
    shift
    ;;
  * )
    errorln "Unknown flag: $key"
    printHelp
    exit 1
    ;;
  esac
  shift
done

if [ "${MODE}" == "up" ]; then
  networkUp
elif [ "${MODE}" == "createChannel" ]; then
  createChannel
elif [ "${MODE}" == "deployCC" ]; then
  deployCC
elif [ "${MODE}" == "start" ]; then
  networkStart
elif [ "${MODE}" == "stop" ]; then
  networkStop
elif [ "${MODE}" == "down" ]; then
  networkDown
else
  printHelp
  exit 1
fi
