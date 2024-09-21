#!/bin/bash

docker run --rm -w "/builder" --platform linux/amd64 -v "${PWD}:/builder" heroiclabs/nakama-pluginbuilder:3.22.0 build -buildvcs=false -buildmode=plugin -trimpath -o ./modules/tictactoe.so
cp modules/tictactoe.so ../nakama/data
