#!/bin/bash

# colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# variables
CURRENT_PATH=`pwd`
CONSUMER_PATH=${GOPATH}/src/github.com/wptechinnovation/wpw-sdk-go/examples/sample-consumer/
PRODUCER_PATH=${GOPATH}/src/github.com/wptechinnovation/wpw-sdk-go/examples/sample-producer-callbacks/
CONSUMER_NAME="sample-consumer"
PRODUCER_NAME="sample-producer-callbacks"
CONSUMER_OUT_NAME="${CONSUMER_NAME}.out"
PRODUCER_OUT_NAME="${PRODUCER_NAME}.out"
CONSUMER_OUT="${CURRENT_PATH}/${CONSUMER_OUT_NAME}"
PRODUCER_OUT="${CURRENT_PATH}/${PRODUCER_OUT_NAME}"

producer_pid=
consumer_pid=

# functions
function usage() {
	echo "Usage: $0 [-c] [-p] [-o <output file>] [-h]"
	echo "Tests the consumer and/or producer"
	echo "Optional options:"

	echo "-o <file>"
	echo "    Puts the output to the file."
	echo " -n no color"
	echo "-h  This help."
	exit 1
}

function build() {
	# build producer
	cd ${PRODUCER_PATH}
	go build
	if [[ $? != 0 ]]; then
		echo -e "${RED}error, producer build failed${NC}"
		exit 1
	fi

	echo -n "."
	# build consumer
	cd ${CONSUMER_PATH}
	go build
	if [[ $? != 0 ]]; then
		echo -e "${RED}error, consumer build failed${NC}"
		exit 1
	fi
	echo -n "."
	echo -e " ${GREEN}ok${NC}"
}

function cleanup() {
	echo -e "${RED}error occured, cleanup, see ${NC}"
	if [[ -n "${consumer_pid}" ]]; then
		# vfy whether consumer is running
		if ps -p ${consumer_pid} > /dev/null 2>&1
		then
			echo "killing consumer"
			kill ${consumer_pid}
		fi
	fi
	
	if [[ -n "${producer_pid}" ]]; then
		# vfy whether producer is running
		if ps -p ${producer_pid} > /dev/null 2>&1
		then
			echo "killing producer"
			kill ${consumer_pid}
		fi
	fi
	exit 1
}

function waitforpid() {
	while ps -p "$1" > /dev/null 2>&1; do
		echo -n "."
		sleep 2
	done
}

trap cleanup SIGINT
trap cleanup SIGTERM

# read startup parameters
while getopts "o:h" o; do
	case "${o}" in
		o)
			o=${OPTARG}
			;;
		n)
			RED=""
			GREEN=""
			NC=""
			;;
		*)
			usage
			;;
	esac
done
shift $((OPTIND-1))

if [[ ${CONSUMER} == 0 && ${PRODUCER} == 0 ]]
then
	# nothing set, so assume we want test both consumer and producer
	CONSUMER=1
	PRODUCER=1
fi

echo "* Initial verification."

echo -n "*** Verify GOPATH is set: "
# GOPATH should be set
if [[ -z ${GOPATH} ]]; then
	echo -e "${RED}error, GOPATH is not set${NC}"
	exit 1
else
	echo -e " ${GOPATH}: ${GREEN}ok${NC}"
fi

# vfy go-lang is installed, consumer, producer paths exists
echo -n "*** Verify Go is installed: "
if ! type go > /dev/null 2>&1
then
	echo -e "${RED}error, cannot find go${NC}"
	exit 1
else
	echo -e "${GREEN}ok${NC}"
fi

# vfy any other consumer / produer is not working
echo -n "*** Verify other producer is not running: "
if ps aux | grep "${PRODUCER_NAME}" | grep -v grep > /dev/null 2>&1
then
	echo -e "${RED}error, other producer is running${NC}"
	exit 1
else
	echo -e "${GREEN}ok${NC}"
fi

echo -n "*** Verify other customer is not running: "
if ps aux | grep "${CONSUMER_NAME}" | grep -v grep > /dev/null 2>&1
then
	echo -e "${RED}error, other customer is running${NC}"
	exit 1
else
	echo -e "${GREEN}ok${NC}"
fi

# vfy that ports are not occupied
echo -n "*** Build consumer and producer: "
build

# start producer in background
cd ${PRODUCER_PATH}
echo "* Start producer."
(./sample-producer-callbacks > "${PRODUCER_OUT}" 2>&1) &
producer_pid=$!

# wait 5 s.
echo -n "*** Wait 5s. for producer to start: "
sleep 5

# vfy that producer is not gone
if ! ps -p ${producer_pid} > /dev/null 2>&1
then
	echo -e "${RED}error, producer gone${NC}"
	exit 1
else
	echo -e "${GREEN}ok${NC}"
fi

# start concumer in background
echo "* Start consumer."
cd ${CONSUMER_PATH}
(./sample-consumer > "${CONSUMER_OUT}" 2>&1) &
consumer_pid=$!


# wait for consumer
echo -n "*** Wait for consumer to finish: "
waitforpid ${consumer_pid}
wait ${consumer_pid}
consumer_status=$?
if [[ ${consumer_status} == 0 ]]; then
	echo -e " ${GREEN}ok${NC}"
else
	echo -e " ${RED}error, consumer exit with ${consumer_status}${NC}"
fi

# producer should be working still
echo -n "*** Verify producer still working: "
if ! ps -p ${producer_pid} > /dev/null 2>&1
then
	echo -e "${RED}error, producer gone${NC}"
	cleanup
else
	echo -e "${GREEN}ok${NC}"
fi

# stop producer
echo "* Stop producer gently."
kill ${producer_pid}

echo -n "*** Verify producer is stopped: "
sleep 3
if ps -p ${producer_pid} > /dev/null 2>&1
then
	echo -e "${RED}error, producer is still alive${NC}"
	cleanup
else
	echo -e "${GREEN}ok${NC}"
fi

echo -e "* ${GREEN}Tests finished successfully.${NC} *"
echo "Note: Outputs can be found in ${PRODUCER_OUT_NAME} and ${CONSUMER_OUT_NAME}"
exit 0
