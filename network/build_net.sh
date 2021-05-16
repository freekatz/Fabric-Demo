#!/bin/bash

bash net.sh up
bash net.sh createChannel
sleep 10
bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep "OR('Org1.member')" -ccm all
bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep "OR('Org2.member')" -ccm all
bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep "AND('Org1.member','Org2.member')" -ccm all
sleep 5
bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep "OR('Org1.member')" -ccm test
bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep "OR('Org2.member')" -ccm test
bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep "AND('Org1.member','Org2.member')" -ccm test
