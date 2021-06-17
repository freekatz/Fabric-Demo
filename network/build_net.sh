#!/bin/bash

. scripts/utils.sh

errorln "Please copy/paste and exec the commands manually."

infoln "bash net.sh up"
infoln ""
infoln "bash net.sh createChannel"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm all"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm all"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep \"OR('Org1.member','Org2.member')\" -ccm all"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 123 -ccn trace -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel123/trace/ -ccep \"OR('Org1.member','Org2.member','Org3.member')\" -ccm all"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 123 -ccn produce -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel123/produce/ -ccep \"OR('Org1.member','Org2.member','Org3.member')\" -ccm all"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 123 -ccn process -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel123/process/ -ccep \"OR('Org1.member','Org2.member','Org3.member')\" -ccm all"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 123 -ccn transport -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel123/transport/ -ccep \"OR('Org1.member','Org2.member','Org3.member')\" -ccm all"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 3 -ccn file -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel3/file/ -ccep \"OR('Org3.member')\" -ccm all"

infoln ""

infoln "bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm test"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm test"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep \"OR('Org1.member','Org2.member')\" -ccm test"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 123 -ccn trace -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel123/trace/ -ccep \"OR('Org1.member','Org2.member','Org3.member')\" -ccm test"
infoln ""
infoln "bash net.sh deployCC -c channel -cid 3 -ccn file -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel3/file/ -ccep \"OR('Org3.member')\" -ccm test"
