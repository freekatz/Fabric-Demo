#!/bin/bash

source scripts/utils.sh

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/orgs/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem


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

println "executing with the following"
println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
println "- CHANNEL_ID: ${C_GREEN}${CHANNEL_ID}${C_RESET}"
println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
println "- CC_SRC_LANGUAGE: ${C_GREEN}${CC_SRC_LANGUAGE}${C_RESET}"
println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
println "- CC_END_POLICY: ${C_GREEN}${CC_END_POLICY}${C_RESET}"
println "- CC_COLL_CONFIG: ${C_GREEN}${CC_COLL_CONFIG}${C_RESET}"
println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"


## 开始

packageChaincode() {
  set -x

  local _P=${PWD}

  cd ${CC_SRC_PATH}
  go env -w GO111MODULE=on
  go env -w GOPROXY=https://goproxy.cn,direct
  go mod vendor

  cd ${_P}
  peer lifecycle chaincode package ${CC_NAME}_${CHANNEL_NAME}${CHANNEL_ID}.tar.gz --path ${CC_SRC_PATH} --label ${CC_NAME}_${CC_VERSION} >&log.txt

  mv ${CC_NAME}_${CHANNEL_NAME}${CHANNEL_ID}.tar.gz ${PWD}/../chaincode/releases/

  docker cp ${PWD}/../chaincode/releases/${CC_NAME}_${CHANNEL_NAME}${CHANNEL_ID}.tar.gz cli10:/opt/gopath/src/github.com/hyperledger/fabric/peer/

  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

# installChaincode PEER ORG
installChaincode() {
  set -x
  mv ${CC_NAME}_${CHANNEL_NAME}${CHANNEL_ID}.tar.gz ${CC_NAME}.tar.gz
  peer lifecycle chaincode install ${CC_NAME}.tar.gz 
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  peer lifecycle chaincode queryinstalled | tee /opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/cc/package_id.txt
}

readPackageID() {
  PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" /opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/cc/package_id.txt)
}

approveForMyOrg() { 
  readPackageID

  set -x
  peer lifecycle chaincode approveformyorg --tls true --cafile ${ORDERER_CA} --channelID ${CHANNEL_NAME}${CHANNEL_ID} --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

commitChaincodeDefinition() {
  set -x
  peer lifecycle chaincode commit  -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA --channelID ${CHANNEL_NAME}${CHANNEL_ID} --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/orgs/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

chaincodeInvokeInit() {
  set -x
  fcn_call='{"function":"'${CC_INIT_FCN}'","Args":[]}'
  infoln "invoke fcn call:${fcn_call}"
  peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C ${CHANNEL_NAME}${CHANNEL_ID} -n ${CC_NAME}  --isInit -c ${fcn_call} --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/orgs/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

chaincodeInvokeTest() {
  ## !!! 输入参数在 shell 中不要有空格
  infoln "登记病人信息"
  chaincodeInvoke '{"function":"registerPatient","Args":["p1","{\"name\":\"ZJH-1\",\"gender\":\"male\",\"birth\":\"1998-10-01\",\"identifyID\":\"xxxxxx-19981001-xxxx\",\"phoneNumber\":\"151-2300-0000\",\"address\":\"ChongQing\",\"nativePlace\":\"NeiMengGu\",\"creditCard\":\"6217-0000-0000-0000\",\"healthcareID\":\"h1\"}"]}' patient
  sleep 5
  infoln "查询登记结果"
  chaincodeInvoke '{"function":"queryPatient","Args":["p1"]}' patient
  sleep 5
  infoln "更新病人信息（名字和性别）"
  chaincodeInvoke '{"function":"updatePatient","Args":["p1","[\"name\",\"gender\"]","[\"ZJH-2\",\"female\"]"]}' patient
  sleep 5
  infoln "查询更新结果"
  chaincodeInvoke '{"function":"queryPatient","Args":["p1"]}' patient
  sleep 5
  infoln "删除病人信息"
  chaincodeInvoke '{"function":"deletePatient","Args":["p1"]}' patient
  sleep 5
  infoln "查询删除结果"
  chaincodeInvoke '{"function":"queryPatient","Args":["p1"]}' patient
  sleep 5
  infoln "再次登记病人信息"
  chaincodeInvoke '{"function":"registerPatient","Args":["p1","{\"name\":\"ZJH-1\",\"gender\":\"male\",\"birth\":\"1998-10-01\",\"identifyID\":\"xxxxxx-19981001-xxxx\",\"phoneNumber\":\"151-2300-0000\",\"address\":\"ChongQing\",\"nativePlace\":\"NeiMengGu\",\"creditCard\":\"6217-0000-0000-0000\",\"healthcareID\":\"h1\"}"]}' patient
  sleep 5
  infoln "查询登记结果"
  chaincodeInvoke '{"function":"queryPatient","Args":["p1"]}' patient
}

chaincodeInvoke() {
  set -x
  local fcn_call=$1
  local cc_name=$2
  infoln "invoke fcn call:${fcn_call}"
  peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA --channelID ${CHANNEL_NAME}${CHANNEL_ID} --name ${cc_name} --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/orgs/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c ${fcn_call}
  >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
}

if [ "${MODE}" == "package" ]; then
  packageChaincode
elif [ "${MODE}" == "install" ]; then
  installChaincode
elif [ "${MODE}" == "approve" ]; then
  approveForMyOrg
elif [ "${MODE}" == "commit" ]; then
  commitChaincodeDefinition
elif [ "${MODE}" == "init" ]; then
  chaincodeInvokeInit
elif [ "${MODE}" == "test" ]; then
  chaincodeInvokeTest
else
  printHelp
  exit 1
fi
