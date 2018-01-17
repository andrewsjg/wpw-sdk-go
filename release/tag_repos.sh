#!/bin/bash

# colors
typeset RED='\033[0;31m'
typeset GREEN='\033[0;32m'
typeset NC='\033[0m'  # No Color

# typeset -r REPO_GOLANG_NAME="wpw-sdk-go"
typeset -r REPO_DOTNET_NAME="wpw-sdk-dotnet"
typeset -r REPO_NODEJS_NAME="wpw-sdk-nodejs"
typeset -r REPO_PYTHON_NAME="wpw-sdk-python"
typeset -r REPO_JAVA_NAME="wpw-sdk-java"
typeset -r REPO_IOT_NAME="wpw-sdk-iot-core"
typeset -r REPO_THRIFT_NAME="wpw-sdk-thrift"
typeset -r REPO_GO_NAME="wpw-sdk-go"
#typeset ALL_REPOS_NAMES="${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} ${REPO_GO_NAME}"

typeset VERSION=

function cleanup {
    echo
}

while true; do
  case "$1" in
    -v | --version ) VERSION="$2"; shift; shift ;;
    -r | --repos_names )
        IN_REPOS_NAMES=(${2//,/ })
        # IFS=','
        # read -ra IN_REPOS_NAMES <<< "$2"
        # #IN_REPOS_NAMES=($2)
        # unset IFS
        shift
        shift
        ;;
    -n | --no-color )
        RED="";
        GREEN="";
        NC="";
        shift ;;
    * ) break ;;
  esac
done

if [[ -z ${VERSION} ]]; then
    echo -e "${RED}error, version name not defined${NC}"
    exit 1
fi

if [[ ${#IN_REPOS_NAMES[@]} -ne 0 ]]; then
    ALL_REPOS_NAMES=("${IN_REPOS_NAMES[@]}")
else
    ALL_REPOS_NAMES=( ${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} ${REPO_GO_NAME} )
fi

echo -e "${GREEN}Tag repos with name: ${VERSION}.${NC}"
# commit repos
for repo_name in ${ALL_REPOS_NAMES[@]};
do
    CURRENT_PATH=`pwd`
    case "${repo_name}" in
        ${REPO_GO_NAME} )
            cd ${WPW_SDK_GO_PATH}/applications/rpc-agent
            RC=$?
            ;;
        ${REPO_IOT_NAME} | ${REPO_THRIFT_NAME} )
            # should be already done with prepare_submodules.sh script
            continue
            ;;
        *) 
            cd ${repo_name}
            RC=$?
            ;;
    esac

    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to change directory to ${repo_name}${NC}"
        cd ${CURRENT_PATH}
        cleanup
        exit 2
    fi

    # git tag -a v0.12-alpha -m "Version 0.12-alpha"
    echo -e "${GREEN}${repo_name}:${NC} git tag -a ${VERSION} -m \"Version ${VERSION}\""
    git tag -a ${VERSION} -m "Version ${VERSION}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to add tag for repo ${repo_name}: git tag -a ${NC}"
        cd ${CURRENT_PATH}
        cleanup
        exit 3
    fi

    cd ${CURRENT_PATH}
done

exit 0
