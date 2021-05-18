#!/bin/bash

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="channel"}

bash scripts/channel/createChannel1.sh $CHANNEL_NAME
bash scripts/channel/createChannel2.sh $CHANNEL_NAME
bash scripts/channel/createChannel3.sh $CHANNEL_NAME
bash scripts/channel/createChannel12.sh $CHANNEL_NAME
bash scripts/channel/createChannel123.sh $CHANNEL_NAME
