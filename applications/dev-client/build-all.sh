#!/bin/bash

echo Cleanup build directory
rm -rf build
mkdir -p build

echo
echo
echo Build dev-client for all supported platforms..
echo
echo
echo "Target LINUX, ARM, GOARM=5"
env GOOS=linux GOARCH=arm GOARM=5 go build && mv dev-client build/dev-client-linux-arm
echo "Target LINUX, 32bit"
env GOOS=linux GOARCH=386 go build && mv dev-client build/dev-client-linux-32
echo "Target LINUX, 64bit"
env GOOS=linux GOARCH=amd64 go build && mv dev-client build/dev-client-linux-64
echo "Target DARWIN, 32bit"
env GOOS=darwin GOARCH=386 go build && mv dev-client build/dev-client-mac-32
echo "Target DARWIN 64bit"
env GOOS=darwin GOARCH=386 go build && mv dev-client build/dev-client-mac-64
echo "Target DARWIN ARM 32bit"
env GOOS=darwin GOARCH=arm go build && mv dev-client build/dev-client-mac-arm
echo "Target WINDOWS, 32bit "
env GOOS=windows GOARCH=386 go build && mv dev-client.exe build/dev-client-win-32.exe
echo "Target WINDOWS, 64bit"
env GOOS=windows GOARCH=amd64 go build && mv dev-client.exe build/dev-client-win-64.exe
