#!/bin/bash

source scripts/utils.sh

CHANNEL_NAME=${1:-"channel"}
CHANNEL_ID=${2:-"1"}
CC_NAME=${3}
CC_SRC_PATH=${4}
CC_SRC_LANGUAGE=${5}
CC_VERSION=${6:-"1.0"}
CC_SEQUENCE=${7:-"1"}
CC_INIT_FCN=${8:-"NA"}
CC_END_POLICY=${9:-"NA"}
CC_COLL_CONFIG=${10:-"NA"}
VERBOSE=${11:-"false"}
MODE=${12:-"package"}

packageChaincode() {
  set -x
  if [ ${CHANNEL_ID} -eq 1 ]; then
    bash ./scripts/cc/deployCC1.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE package
  elif [ ${CHANNEL_ID} -eq 2 ]; then
    bash ./scripts/cc/deployCC2.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE package
  elif [ ${CHANNEL_ID} -eq 3 ]; then
    bash ./scripts/cc/deployCC3.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE package
  elif [ ${CHANNEL_ID} -eq 12 ]; then
    bash ./scripts/cc/deployCC12.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE package
  elif [ ${CHANNEL_ID} -eq 123 ]; then
    bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE package
  else
    errorln "Channel not exist!"
    exit 1
  fi
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

# installChaincode PEER ORG
installChaincode() {
  set -x
  if [ ${CHANNEL_ID} -eq 1 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC1.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install
  elif [ ${CHANNEL_ID} -eq 2 ]; then
    docker exec cli20 bash ./scripts/cc/deployCC2.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install
  elif [ ${CHANNEL_ID} -eq 3 ]; then
    docker exec cli30 bash ./scripts/cc/deployCC3.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install
  elif [ ${CHANNEL_ID} -eq 12 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC12.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install

    docker exec cli20 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install
  elif [ ${CHANNEL_ID} -eq 123 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install

    docker exec cli20 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install

    docker exec cli30 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE install

  else
    errorln "Channel not exist!"
    exit 1
  fi
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

# approveForMyOrg VERSION PEER ORG
approveForMyOrg() {
  set -x
  if [ ${CHANNEL_ID} -eq 1 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC1.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve
  elif [ ${CHANNEL_ID} -eq 2 ]; then
    docker exec cli20 bash ./scripts/cc/deployCC2.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve
  elif [ ${CHANNEL_ID} -eq 3 ]; then
    docker exec cli30 bash ./scripts/cc/deployCC3.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve
  elif [ ${CHANNEL_ID} -eq 12 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC12.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve

    docker exec cli20 bash ./scripts/cc/deployCC12.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve
  elif [ ${CHANNEL_ID} -eq 123 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve

    docker exec cli20 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve

    docker exec cli30 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE approve

  else
    errorln "Channel not exist!"
    exit 1
  fi
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

# commitChaincodeDefinition VERSION PEER ORG (PEER ORG)...
commitChaincodeDefinition() {
  set -x
  if [ ${CHANNEL_ID} -eq 1 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC1.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE commit
  elif [ ${CHANNEL_ID} -eq 2 ]; then
    docker exec cli20 bash ./scripts/cc/deployCC2.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE commit
  elif [ ${CHANNEL_ID} -eq 3 ]; then
    docker exec cli30 bash ./scripts/cc/deployCC3.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE commit
  elif [ ${CHANNEL_ID} -eq 12 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC12.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE commit
  elif [ ${CHANNEL_ID} -eq 123 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE commit
  else
    errorln "Channel not exist!"
    exit 1
  fi
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}


chaincodeInvokeInit() {
  set -x
  if [ ${CHANNEL_ID} -eq 1 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC1.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE init
  elif [ ${CHANNEL_ID} -eq 2 ]; then
    docker exec cli20 bash ./scripts/cc/deployCC2.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE init
  elif [ ${CHANNEL_ID} -eq 3 ]; then
    docker exec cli30 bash ./scripts/cc/deployCC3.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE init
  elif [ ${CHANNEL_ID} -eq 12 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC12.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE init
  elif [ ${CHANNEL_ID} -eq 123 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE init
  else
    errorln "Channel not exist!"
    exit 1
  fi
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

chaincodeInvokeTest() {
  set -x
  if [ ${CHANNEL_ID} -eq 1 ]; then
    docker exec cli10 bash ./scripts/cc/deployCC1.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE test
  elif [ ${CHANNEL_ID} -eq 2 ]; then
    docker exec cli20 bash ./scripts/cc/deployCC2.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE test
  elif [ ${CHANNEL_ID} -eq 3 ]; then
    docker exec cli30 bash ./scripts/cc/deployCC3.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE test
  elif [ ${CHANNEL_ID} -eq 12 ]; then
    # todo 修改 msp id
    docker exec cli10 bash ./scripts/cc/deployCC12.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE test
  elif [ ${CHANNEL_ID} -eq 123 ]; then
    # todo 修改 msp id
    docker exec cli10 bash ./scripts/cc/deployCC123.sh $CHANNEL_NAME $CHANNEL_ID $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $VERBOSE test
  else
    errorln "Channel not exist!"
    exit 1
  fi
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

if [ "${MODE}" == "all" ]; then
  packageChaincode
  installChaincode
  sleep 3
  approveForMyOrg
  sleep 3
  commitChaincodeDefinition
  sleep 3
  if [ "$CC_INIT_FCN" = "NA" ]; then
    infoln "Chaincode initialization is not required"
  else
    chaincodeInvokeInit
  fi
elif [ "${MODE}" == "package" ]; then
  packageChaincode
elif [ "${MODE}" == "install" ]; then
  installChaincode
elif [ "${MODE}" == "approve" ]; then
  approveForMyOrg
elif [ "${MODE}" == "commit" ]; then
  commitChaincodeDefinition
elif [ "${MODE}" == "init" ]; then
  if [ "$CC_INIT_FCN" = "NA" ]; then
    infoln "Chaincode initialization is not required"
  else
    chaincodeInvokeInit
  fi
elif [ "${MODE}" == "test" ]; then
  chaincodeInvokeTest
else
  printHelp
  exit 1
fi