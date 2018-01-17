#!/bin/bash

# colors
typeset RED='\033[0;31m'
typeset GREEN='\033[0;32m'
typeset NC='\033[0m'  # No Color

typeset -r REPO_GO_NAME="wpw-sdk-go"
typeset -r REPO_IOT_NAME="wpw-sdk-iot-core"
typeset -r REPO_THRIFT_NAME="wpw-sdk-thrift"

typeset RC_BRANCH_NAME=""
typeset RC_MASTER_BRANCH_NAME=""

typeset WPW_SDK_GO_PATH=${GOPATH}/src/github.com/WPTechInnovation/${REPO_GO_NAME}
typeset VERSION=""
typeset ADD_TAG=""

START_DIR=`pwd`

function cleanup {
    echo
}



while true; do
  case "$1" in
    -v | --version )
        VERSION="$2"
        shift
        shift
        ;;
    -t | --add_tag )
        ADD_TAG="-t"
        shift
        ;;
    -b | --branch )
        RC_BRANCH_NAME="$2";
        shift
        shift
        ;;
    -m | --master_branch )
        RC_MASTER_BRANCH_NAME="$2"
        shift
        shift
        ;;
    -r | --repos_names )
        IN_REPOS_NAMES=(${2//,/ })
        shift
        shift
        ;;
    -n | --no-color )
        RED="";
        GREEN="";
        NC="";
        shift
        ;;
    * ) break ;;
  esac
done

# vfy input parameters
if [[ -z ${VERSION} ]]; then
    echo -e "${RED}error, version name not defined${NC}"
    exit 1
fi

if [[ -z ${RC_MASTER_BRANCH_NAME} ]]; then
    echo -e "${RED}error, master branch name not defined${NC}"
    exit 1
fi

# prepare env values
if [[ ${#IN_REPOS_NAMES[@]} -ne 0 ]]; then
    ALL_REPOS_NAMES=("${IN_REPOS_NAMES[@]}")
else
    ALL_REPOS_NAMES=( ${REPO_GO_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} )
fi

# vfy that version not exists
if [[ $(git tag -l "${VERSION}") ]]; then
    echo -e "${RED}error, the version/tag ${VERSION} is already defined"
    exit 2
fi


# build_rpc_agents 

START_PATH=`pwd`
cd ${WPW_SDK_GO_PATH}/applications/rpc-agent
RC=$?
if [[ ${RC} != 0 ]]
then
    echo -e "${RED}error, failed to change directory to ${WPW_SDK_GO_PATH}/applications/rpc-agent${NC}"
    cd ${START_PATH}
    cleanup
    exit 2
fi

./build-all.sh -v ${VERSION}
RC=$?
if [[ ${RC} != 0 ]]
then
    echo -e "${RED}error, failed to build RPC agents${NC}"
    cd ${START_PATH}
    cleanup
    exit 2
fi

cd ${START_PATH}

# copy rpc agents to iot directory
cp ${WPW_SDK_GO_PATH}/applications/rpc-agent/build/rpc* ${REPO_IOT_NAME}/bin/
RC=$?
if [[ ${RC} != 0 ]]
then
    echo -e "${RED}error, failed to copy rpc agents to ${REPO_IOT_NAME}/bin/${NC}"
    cleanup
    exit 2
fi

# commit_iot
repo_name=${REPO_IOT_NAME}
cd ${repo_name}
RC=$?
if [[ ${RC} != 0 ]]
then
    echo -e "${RED}error, failed to change directory to ${repo_name}${NC}"
    cd ${START_PATH}
    cleanup
    exit 2
fi

# vfy if branch name is correct
CURRENT_BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
if [[ ${CURRENT_BRANCH_NAME} != "${RC_MASTER_BRANCH_NAME}" ]]; then
    echo -e "${RED}error, current branch name ${CURRENT_BRANCH_NAME} is different than ${RC_MASTER_BRANCH_NAME} for ${repo_name}${NC}"
    cd ${START_DIR}
    cleanup
    exit 1
fi

echo -e "${GREEN}${repo_name}:${NC} git commit -a -m \"new rpc-agents for version ${VERSION}\""
git commit -a -m "new rpc-agents for version ${VERSION}"
RC=$?
if [[ ${RC} != 0 ]]
then
    echo -e "${RED}error, failed to commit changes for new rpc-agents version ${VERSION}${NC}"
    cd ${START_PATH}
    cleanup
    exit 2
fi

# tag_repo
if [[ -n "${ADD_TAG}" ]]; then
    echo -e "${GREEN}${repo_name}:${NC} git tag -a ${VERSION} -m \"Version ${VERSION}\""
    git tag -a ${VERSION} -m "Version ${VERSION}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to add tag for repo ${repo_name}: git tag -a ${VERSION} -m \"Version ${VERSION}\"${NC}"
        cd ${START_DIR}
        cleanup
        exit 3
    fi
fi

# push changes
echo -e "${GREEN}Push repo ${repo_name}${NC}"
echo -e "${GREEN}${repo_name}:${NC} git push origin ${RC_MASTER_BRANCH_NAME}"
git push origin ${RC_MASTER_BRANCH_NAME}
RC=$?
if [[ ${RC} != 0 ]]
then
    echo -e "${RED}error, failed to: git push origin ${RC_MASTER_BRANCH_NAME}${NC}"
    cd ${START_DIR}
    cleanup
    exit 4
fi

if [[ -n "${ADD_TAG}" ]]; then
    # push tags
    echo -e "${GREEN}${repo_name}:${NC} git push --tags"
    git push --tags
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to: \"git push --tags\" in ${repo_name}${NC}"
        cd ${START_DIR}
        cleanup
        exit 5
    fi
fi

cd ${START_DIR}

# commit thrift
# not finished yet, exit 0 for now
# Thrift is specific and contain master only officially,
# continue for master branch only
if [[ "${RC_MASTER_BRANCH_NAME}" == "master" ]]; then

    repo_name=${REPO_THRIFT_NAME}
    cd ${repo_name}
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to change directory to ${repo_name}${NC}"
        cd ${START_PATH}
        cleanup
        exit 2
    fi

    echo -e "${GREEN}${repo_name}:${NC} git checkout ${RC_MASTER_BRANCH_NAME}"
    git checkout "${RC_MASTER_BRANCH_NAME}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to checkout to ${RC_MASTER_BRANCH_NAME} for repo ${repo_name}"
        cd ${START_PATH}
        cleanup
        exit 3
    fi

    echo -e "${GREEN}${repo_name}:${NC} git pull"
    git pull
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to: git pull${NC}"
        cd ${START_PATH}
        cleanup
        exit 4
    fi

    # echo -e "${GREEN}${repo_name}:${NC} git merge --no-ff --no-commit ${RC_BRANCH_NAME}"
    # git merge --no-ff --no-commit "${RC_BRANCH_NAME}"
    echo -e "${GREEN}${repo_name}:${NC} git merge --no-ff --no-edit ${RC_BRANCH_NAME}"
    git merge --no-ff --no-edit "${RC_BRANCH_NAME}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to merge ${RC_BRANCH_NAME} to ${RC_MASTER_BRANCH_NAME}${NC}"
        cd ${START_PATH}
        cleanup
        exit 5
    fi

    # push changes
    echo -e "${GREEN}Push repo ${repo_name}${NC}"
    echo -e "${GREEN}${repo_name}:${NC} git push origin ${RC_MASTER_BRANCH_NAME}"
    git push origin ${RC_MASTER_BRANCH_NAME}
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to: git push origin ${RC_MASTER_BRANCH_NAME}${NC}"
        cd ${START_DIR}
        cleanup
        exit 4
    fi

    # push tags
    if [[ -n "${ADD_TAG}" ]]; then
        echo -e "${GREEN}${repo_name}:${NC} git push --tags"
        git push --tags
        RC=$?
        if [[ ${RC} != 0 ]]; then
            echo -e "${RED}error, failed to: git push --tags in ${repo_name}${NC}"
            cd ${START_DIR}
            cleanup
            exit 5
        fi
    fi

    # go back to initial directory
    cd ${START_DIR}

fi

exit 0
