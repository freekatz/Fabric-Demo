#!/bin/bash

. scripts/utils.sh

errorln "Please do not exec this script directly, and use cat build_net.sh to copy/paste the commands."

infoln "bash net.sh up"

infoln "bash net.sh createChannel"

infoln "bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 1 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/<>/ -ccep \"OR('Org1.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 2 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/<>/ -ccep \"OR('Org2.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 3 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel3/<>/ -ccep \"OR('Org3.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep \"OR('Org1.member','Org2.member')\" -ccm all"
infoln "bash net.sh deployCC -c channel -cid 123 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel123/<>/ -ccep \"OR('Org1.member','Org2.member','Org3.member')\" -ccm all"

infoln "bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 1 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/<>/ -ccep \"OR('Org1.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 2 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/<>/ -ccep \"OR('Org2.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 3 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel3/<>/ -ccep \"OR('Org3.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep \"OR('Org1.member','Org2.member')\" -ccm test"
infoln "bash net.sh deployCC -c channel -cid 123 -ccn <> -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel123/<>/ -ccep \"OR('Org1.member','Org2.member','Org3.member')\" -ccm test"