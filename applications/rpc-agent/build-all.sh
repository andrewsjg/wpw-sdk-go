#!/bin/bash

while getopts ":v:" opt; do
  case $opt in
    v) version="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done

if [[ -z "${version}" ]] ; then
  echo "-v (build version) is required"
  exit
fi

echo Cleanup build directory
rm -rf build
mkdir -p build

echo
echo
echo "Build rpc-agent for all supported platforms.."
echo
echo "Will set version = $version and date = `date -u +%d-%m-%Y@%H:%M:%S`"
echo
echo "Target LINUX, 32bit"
env GOOS=linux GOARCH=386 go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=LINUX_386" -o build/rpc-agent-linux-386 main.go
echo "Target LINUX, 64bit"
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=LINUX_AMD64" -o build/rpc-agent-linux-amd64 main.go
echo "Target LINUX, ARM, GOARM=5"
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=LINUX_ARM" -o build/rpc-agent-linux-arm32 main.go
echo "Target LINUX, ARM64, GOARM=5"
env GOOS=linux GOARCH=arm64 GOARM=5 go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=LINUX_ARM64" -o build/rpc-agent-linux-arm64 main.go
echo "Target DARWIN, 32bit"
env GOOS=darwin GOARCH=386 go build && go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=DARWIN_386" -o build/rpc-agent-darwin-386 main.go
echo "Target DARWIN, 64bit"
env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=DARWIN_AMD64" -o build/rpc-agent-darwin-amd64 main.go
echo "Target WINDOWS, 32bit "
env GOOS=windows GOARCH=386 go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=WIN_386" -o build/rpc-agent-windows-386.exe main.go
echo "Target WINDOWS, 64bit"
env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.applicationVersion=$version -X main.applicationBuildDate=`date -u +%d-%m-%Y@%H:%M:%S` -X main.applicationPlatform=WIN_AMD64" -o build/rpc-agent-windows-amd64.exe main.go

#cd build && zip -r -X rpc-agent-bins-$version.zip * && cd ../
#cd build && tar czf rpc-agent-bins-$version.tar.gz * && cd ../
