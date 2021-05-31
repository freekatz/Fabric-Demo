
#!/bin/bash

export CLUSTER_SECRET=b8def8a4ab09bdbb6bc3806115f6437184297f94c703c3729360fb11d0d89348

. utils.sh

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

function networkUp() {
    docker-compose up -d
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
  infoln "Down the ipfs"
  docker-compose down --volumes --remove-orphans
  # Bring down the network, deleting the volumes
  #Cleanup the chaincode containers
  clearContainers
  #Cleanup images
  removeUnwantedImages

  cd $_p
  res=$?
  { set +x; } 2>/dev/null
}

## Parse mode
if [[ $# -lt 1 ]] ; then
  printHelp
  exit 0
else
  MODE=$1
  shift
fi

if [ "${MODE}" == "up" ]; then
  networkUp
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
