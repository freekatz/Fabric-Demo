#!/bin/bash

. scripts/utils.sh

errorln "Please do not exec this script directly, and use cat build_net.sh to copy/paste the commands."

infoln "bash net.sh up"

infoln "bash net.sh createChannel"

infoln "## 这里设置的背书策略要与 SDK 使用的一致，否则有时会出现异常"
infoln "bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccm all"

infoln "bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccm test"
