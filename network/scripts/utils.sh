#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# Print the usage message
function printHelp() {
  USAGE="$1"
  if [ "$USAGE" == "up" ]; then
    println "Usage: "
    println "  net.sh \033[0;32mup\033[0m [Flags]"
    println
    println "    Flags:"
    println "    -verbose - Verbose mode"
    println
    println "    -h - Print this message"
    println
    println " Examples:"
    println "   bash net.sh up"
  elif [ "$USAGE" == "createChannel" ]; then
    println "Usage: "
    println "  net.sh \033[0;32mcreateChannel\033[0m [Flags]"
    println
    println "    Flags:"
    println "    -c <channel base name> - Name of channel to create (defaults to \"channel\") note that: the channel id has been defined!"
    println "    -verbose - Verbose mode"
    println
    println "    -h - Print this message"
    println
    println " Examples:"
    println "   bash net.sh createChannel -c channel"
  elif [ "$USAGE" == "deployCC" ]; then
    println "Usage: "
    println "  net.sh \033[0;32mdeployCC\033[0m [Flags]"
    println
    println "    Flags:"
    println "    -c <channel name> - Name of channel to deploy chaincode to"
    println "    -cid <channel id> - ID of channel to deploy chaincode to"
    println "    -ccn <name> - Chaincode name."
    println "    -ccl <language> - Programming language of chaincode to deploy: go, java, javascript, typescript"
    println "    -ccv <version>  - Chaincode version. 1.0 (default), v2, version3.x, etc"
    println "    -ccs <sequence>  - Chaincode definition sequence. Must be an integer, 1 (default), 2, 3, etc"
    println "    -ccp <path>  - File path to the chaincode."
    println "    -ccep <policy>  - (Optional) Chaincode endorsement policy using signature policy syntax. The default policy requires an endorsement from Orgs"
    println "    -cccg <collection-config>  - (Optional) File path to private data collections configuration file"
    println "    -cci <fcn name>  - (Optional) Name of chaincode initialization function. When a function is provided, the execution of init will be requested and the function will be invoked."
    println "    -ccm <chaincode mode>  - (Optional) included: package, install, approve, commit, init, test, all."
    println
    println "    -h - Print this message"
    println
    println " Examples:"
    println "   bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm all"
    println "   bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccm test"
    println "   bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm all"
    println "   bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep \"AND('Org1.member','Org2.member')\" -ccm all"
  else
    println "Usage: "
    println "  net.sh <Mode> [Flags]"
    println "    Modes:"
    println "      \033[0;32mup\033[0m - Bring up Fabric orderer and peer nodes. No channel is created"
    println "      \033[0;32mcreateChannel\033[0m - Create and join a channel after the net is created"
    println "      \033[0;32mdeployCC\033[0m - Deploy a chaincode to a channel (defaults to asset-transfer-basic)"
    println "      \033[0;32mstart\033[0m - Start the net"
    println "      \033[0;32mstop\033[0m - Stop the net"
    println "      \033[0;32mrestart\033[0m - Restart the net"
    println "      \033[0;32mdown\033[0m - Bring down the net"
    println
    println "    Flags:"
    println "    Used with \033[0;32mnet.sh up\033[0m, \033[0;32mnet.sh createChannel\033[0m:"
    println "    -c <channel base name> - Name of channel to create (defaults to \"channel\") note that: the channel id has been defined!"
    println "    -verbose - Verbose mode"
    println
    println "    Used with \033[0;32mnet.sh deployCC\033[0m"
    println "    -c <channel name> - Name of channel to deploy chaincode to"
    println "    -cid <channel id> - ID of channel to deploy chaincode to"
    println "    -ccn <name> - Chaincode name."
    println "    -ccl <language> - Programming language of the chaincode to deploy: go, java, javascript, typescript"
    println "    -ccv <version>  - Chaincode version. 1.0 (default), v2, version3.x, etc"
    println "    -ccs <sequence>  - Chaincode definition sequence. Must be an integer, 1 (default), 2, 3, etc"
    println "    -ccp <path>  - File path to the chaincode."
    println "    -ccep <policy>  - (Optional) Chaincode endorsement policy using signature policy syntax. The default policy requires an endorsement from Orgs"
    println "    -cccg <collection-config>  - (Optional) File path to private data collections configuration file"
    println "    -cci <fcn name>  - (Optional) Name of chaincode initialization function. When a function is provided, the execution of init will be requested and the function will be invoked."
    println "    -ccm <chaincode mode>  - (Optional) included: package, install, approve, commit, init, test, all."
    println
    println "    -h - Print this message"
    println
    println " Examples:"
    println "   bash net.sh up"
    println "   bash net.sh createChannel -c channel"
    println "   bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccep \"OR('Org1.member')\" -ccm all"
    println "   bash net.sh deployCC -c channel -cid 1 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel1/patient/ -ccm test"
    println "   bash net.sh deployCC -c channel -cid 2 -ccn patient -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel2/patient/ -ccep \"OR('Org2.member')\" -ccm all"
    println "   bash net.sh deployCC -c channel -cid 12 -ccn bridge -ccl go -ccv 1.0 -ccs 1 -ccp ../chaincode/channel12/bridge/ -ccep \"AND('Org1.member','Org2.member')\" -ccm all"
  fi
}

# println echos string
function println() {
  echo -e "$1"
}

# errorln echos i red color
function errorln() {
  println "${C_RED}${1}${C_RESET}"
}

# successln echos in green color
function successln() {
  println "${C_GREEN}${1}${C_RESET}"
}

# infoln echos in blue color
function infoln() {
  println "${C_BLUE}${1}${C_RESET}"
}

# warnln echos in yellow color
function warnln() {
  println "${C_YELLOW}${1}${C_RESET}"
}

# fatalln echos in red color and exits with fail status
function fatalln() {
  errorln "$1"
  exit 1
}

export -f errorln
export -f successln
export -f infoln
export -f warnln
