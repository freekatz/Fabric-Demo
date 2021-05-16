#!/bin/bash

CHANNEL_NAME="$1"
CHANNEL_ID="$2"
: ${CHANNEL_NAME:="channel"}
: ${CHANNEL_ID:="1"}

peer channel join -b ${CHANNEL_NAME}${CHANNEL_ID}.block
