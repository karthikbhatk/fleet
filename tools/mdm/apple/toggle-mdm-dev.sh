#!/bin/bash

# To toggle MDM, run `source toggle-mdm-dev`

if [[ $USE_MDM == "1" ]]; then
export USE_MDM=0
else
export USE_MDM=1
fi

source $FLEET_ENV_PATH
