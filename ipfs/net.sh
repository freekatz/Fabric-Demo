
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

function networkInit() {
  set -x
  infoln "允许跨域访问 API"
  docker exec ipfs0 ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]'
  docker exec ipfs0 ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'

  docker exec ipfs1 ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]'
  docker exec ipfs1 ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'

  infoln "分发 swarm key"
  go get -u github.com/Kubuxu/go-ipfs-swarm-key-gen/ipfs-swarm-key-gen
  if [ -f "~/.ipfs/swarm.key" ]; then
    rm -Rf ipfs-swarm-key-gen > ~/.ipfs/swarm.key
  fi
  ipfs-swarm-key-gen > ~/.ipfs/swarm.key

  docker cp ~/.ipfs/swarm.key ipfs0:/data/
  docker cp ~/.ipfs/swarm.key ipfs1:/data/

  infoln "移除其他全球节点"
  docker exec ipfs0 ipfs bootstrap rm --all
  docker exec ipfs1 ipfs bootstrap rm --all

  infoln "# 暂时还需要手动配置"
  infoln "# ipfs1 添加节点 ipfs0"
  infoln "docker exec ipfs0 ipfs id"
  infoln "docker exec ipfs1 ipfs bootstrap add <address>"
  # docker exec ipfs1 ipfs bootstrap add /ip4/172.20.0.3/tcp/4001/ipfs/QmWGWQovWv9jBdxvtPNoCjuu5f8gY5ovTi1Ln2fDgMSwfs

  infoln "# ipfs1 添加节点 ipfs0"
  infoln "docker exec ipfs1 ipfs id"
  infoln "docker exec ipfs0 ipfs bootstrap add <address>"
  # docker exec ipfs0 ipfs bootstrap add /ip4/172.20.0.2/tcp/4001/ipfs/QmU8r2eekcjsoeGQoKypKaDedFDR2UT9PC8XisBzgWC6yT

  infoln "# 重启"
  infoln "bash net.sh restart"

  infoln "# 测试"
  infoln "sleep 3"
  infoln "docker exec ipfs0 ipfs swarm peers"

  res=$?
  { set +x; } 2>/dev/null
}

function networkUp() {
    docker-compose up -d
}

function networkStart() {
  docker-compose start
}

function networkStop() {
  docker-compose stop
}

function networkRestart() {
  docker-compose restart
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
elif [ "${MODE}" == "init" ]; then
  networkInit
elif [ "${MODE}" == "start" ]; then
  networkStart
elif [ "${MODE}" == "stop" ]; then
  networkStop
elif [ "${MODE}" == "restart" ]; then
  networkRestart
elif [ "${MODE}" == "down" ]; then
  networkDown
else
  printHelp
  exit 1
fi
