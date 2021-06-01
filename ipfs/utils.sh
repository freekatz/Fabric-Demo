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
    println "    -h - Print this message"
    println
    println " Examples:"
    println "   bash net.sh up"
  else
    println "Usage: "
    println "  net.sh <Mode> [Flags]"
    println "    Modes:"
    println "      \033[0;32mup\033[0m - Bring up IPFS nodes"
    println "      \033[0;32minit\033[0m - Initialize the IPFS network"
    println "      \033[0;32mstart\033[0m - Start the net"
    println "      \033[0;32mstop\033[0m - Stop the net"
    println "      \033[0;32mrestart\033[0m - Restart the net"
    println "      \033[0;32mdown\033[0m - Bring down the net"
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
