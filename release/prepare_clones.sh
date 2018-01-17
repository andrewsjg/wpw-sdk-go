#!/bin/bash

# colors
typeset RED='\033[0;31m'
typeset GREEN='\033[0;32m'
typeset NC='\033[0m'  # No Color

# setup default values if empty
[[ -z "${MIRRORS}" ]] &&        MIRRORS="https://github.com/WPTechInnovation"

[[ -z "${REPO_GO}" ]] &&        REPO_GO="${MIRRORS}/wpw-sdk-go.git"
[[ -z "${REPO_DOTNET}" ]] &&    REPO_DOTNET="${MIRRORS}/wpw-sdk-dotnet.git"
[[ -z "${REPO_NODEJS}" ]] &&    REPO_NODEJS="${MIRRORS}/wpw-sdk-nodejs.git"
[[ -z "${REPO_PYTHON}" ]] &&    REPO_PYTHON="${MIRRORS}/wpw-sdk-python.git"
[[ -z "${REPO_JAVA}" ]] &&      REPO_JAVA="${MIRRORS}/wpw-sdk-java.git"
[[ -z "${REPO_IOT}" ]] &&       REPO_IOT="${MIRRORS}/wpw-sdk-iot-core.git"
[[ -z "${REPO_THRIFT}" ]] &&    REPO_THRIFT="${MIRRORS}/wpw-sdk-thrift.git"
#typeset ALL_REPOS=( ${REPO_GO} ${REPO_DOTNET} ${REPO_NODEJS} ${REPO_PYTHON} ${REPO_JAVA} ${REPO_IOT} ${REPO_THRIFT} )

typeset -r REPO_GO_NAME="wpw-sdk-go"
typeset -r REPO_DOTNET_NAME="wpw-sdk-dotnet"
typeset -r REPO_NODEJS_NAME="wpw-sdk-nodejs"
typeset -r REPO_PYTHON_NAME="wpw-sdk-python"
typeset -r REPO_JAVA_NAME="wpw-sdk-java"
typeset -r REPO_IOT_NAME="wpw-sdk-iot-core"
typeset -r REPO_THRIFT_NAME="wpw-sdk-thrift"
typeset ALL_REPOS_NAMES=()

typeset RC_BRANCH_NAME=""

CURRENT_PATH=`pwd`
export GOPATH=${CURRENT_PATH}/go
export WPW_SDK_GO_PATH=${GOPATH}/src/github.com/WPTechInnovation/${REPO_GO_NAME}

function cleanup {
    echo
}

while true; do
  case "$1" in
    -b | --branch ) RC_BRANCH_NAME="$2";
        shift
        shift
        ;;
    -r | --repos_names )
        IN_REPOS_NAMES=(${2//,/ })
        # IFS=','
        # read -ra IN_REPOS_NAMES <<< "$2"
        # #IN_REPOS_NAMES=($2)
        # unset IFS
        shift
        shift
        ;;
    -e | --repos )
        IFS=','
        IN_REPOS=($2)
        unset IFS
        shift
        shift
        ;;
    -n | --no-color )
        RED="";
        GREEN="";
        NC="";
        shift
        ;;
    * )
        # echo -e "${RED}warning, unexpected input parameter {$1}${NC}"
        # exit 1
        break
        ;;
  esac
done

if [[ -z ${RC_BRANCH_NAME} ]]; then
    echo -e "${RED}error, branch name not defined${NC}"
    exit 1
fi

# check if 
if [[ ${#IN_REPOS_NAMES[@]} -ne 0 ]]; then
    ALL_REPOS_NAMES=("${IN_REPOS_NAMES[@]}")
else
    ALL_REPOS_NAMES=( ${REPO_GO_NAME} ${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} )
fi

if [[ ${#IN_REPOS[@]} -ne 0 ]]; then
    ALL_REPOS=("${IN_REPOS[@]}")
else
    ALL_REPOS=( ${REPO_GO} ${REPO_DOTNET} ${REPO_NODEJS} ${REPO_PYTHON} ${REPO_JAVA} ${REPO_IOT} ${REPO_THRIFT} )
fi

function prepareGoEnv {
    local CURRENT_PATH=`pwd`
    cd ${WPW_SDK_GO_PATH}/applications/rpc-agent
    echo -e "${GREEN}${repo_name}:${NC} git checkout ${RC_BRANCH_NAME}"
    git checkout develop
    # go get without building
    go get -d
    cd ../../../../../git.apache.org/thrift.git/
    echo -e "${GREEN}${repo_name}:${NC} changing the thrift to version 0.10.0"
    git checkout 0.10.0
    cd ${CURRENT_PATH}
}

echo -e "${GREEN}Cloning all repos.${NC}"

# clone repos
for repo in ${ALL_REPOS[@]};
do
    case "${repo}" in
        ${REPO_GO} )
            # go repo is cloned in a different way than others
            mkdir -p ${GOPATH}/src/github.com/WPTechInnovation/
            CURRENT_PATH=`pwd`
            cd ${GOPATH}/src/github.com/WPTechInnovation/
            git clone ${repo}
            RC=$?
            if [[ ${RC} != 0 ]]
            then
                echo -e "${RED}error, failed to clone ${repo}${NC}"
                cd ${CURRENT_PATH}
                cleanup
                exit 2
            fi
            cd ${CURRENT_PATH}
            ;;
        * )
            echo -e "${GREEN}git clone ${repo}${NC}"
            git clone ${repo}
            RC=$?
            if [[ ${RC} != 0 ]]
            then
                echo -e "${RED}error, failed to clone ${repo}${NC}"
                cleanup
                exit 2
            fi
            ;;
    esac
done

echo -e "${GREEN}Changing the branch to ${RC_BRANCH_NAME}.${NC}"
# change branch
for repo_name in ${ALL_REPOS_NAMES[@]};
do
    case "${repo_name}" in
        ${REPO_IOT_NAME} )
            # ignore changing branch to develop for wpw-sdk-iot-core
            continue
            ;;
        ${REPO_THRIFT_NAME} )
            continue
            ;;
        ${REPO_GO_NAME} )
            # already done in a loop above
            prepareGoEnv
            continue
            ;;
        *)
            ;;
    esac

    cd ${repo_name}
    echo -e "${GREEN}${repo_name}:${NC} git checkout ${RC_BRANCH_NAME}"
    git checkout "${RC_BRANCH_NAME}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to checkout ${repo_name} to ${RC_BRANCH_NAME}${NC}"
        cd ..
        cleanup
        exit 3
    fi
    cd ..
done

exit 0
