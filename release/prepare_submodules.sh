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

function die {
    echo -e "${1}" >&2
    cleanup
    exit 1
}

function usage {
    echo "usage: prepare_submodules.sh -v|--version <version>"
    echo "                             -b|--branch <source branch>"
    echo "                             -m|--master_branch <destination branch>"
    echo "                             [-t|--add_tag]"
    echo "                             [-r|--repos_names <coma separated repos name>]"
    echo "                             [-n|--no-color]"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -v | --version )
        VERSION="$2"
        shift
        ;;
    -t | --add_tag )
        ADD_TAG="-t"
        ;;
    -b | --branch )
        RC_BRANCH_NAME="$2";
        shift
        ;;
    -m | --master_branch )
        RC_MASTER_BRANCH_NAME="$2"
        shift
        ;;
    -r | --repos_names )
        IN_REPOS_NAMES=(${2//,/ })
        shift
        ;;
    -n | --no-color )
        RED="";
        GREEN="";
        NC="";
        ;;
    * )
        echo -e "${RED}Invalid argument ${1}${NC}" >&2
        usage
        exit 1
        ;;
  esac
  shift
done

# vfy input parameters
if [[ -z ${VERSION} ]]; then
    die "${RED}error, version name not defined${NC}"
fi

if [[ -z ${RC_MASTER_BRANCH_NAME} ]]; then
    die "${RED}error, master branch name not defined${NC}"
fi

# prepare env values
if [[ ${#IN_REPOS_NAMES[@]} -ne 0 ]]; then
    ALL_REPOS_NAMES=("${IN_REPOS_NAMES[@]}")
else
    ALL_REPOS_NAMES=( ${REPO_GO_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} )
fi

# vfy that version not exists
if [[ $(git tag -l "${VERSION}") ]]; then
    die "${RED}error, the version/tag ${VERSION} is already defined"
fi


# build_rpc_agents 

START_PATH=`pwd`
cd ${WPW_SDK_GO_PATH}/applications/rpc-agent
RC=$?
if [[ ${RC} -ne 0 ]]
then
    cd "${START_PATH}"
    die "${RED}error, failed to change directory to ${WPW_SDK_GO_PATH}/applications/rpc-agent${NC}"
fi

./build-all.sh -v ${VERSION}
RC=$?
if [[ ${RC} -ne 0 ]]; then
    cd "${START_PATH}"
    die "${RED}error, failed to build RPC agents${NC}"
fi

cd ${START_PATH}

# copy rpc agents to iot directory
cp ${WPW_SDK_GO_PATH}/applications/rpc-agent/build/rpc* ${REPO_IOT_NAME}/bin/
RC=$?
if [[ ${RC} -ne 0 ]]; then
    die "${RED}error, failed to copy rpc agents to ${REPO_IOT_NAME}/bin/${NC}"
fi

# commit_iot
repo_name=${REPO_IOT_NAME}
cd ${repo_name}
RC=$?
if [[ ${RC} -ne 0 ]]; then
    cd "${START_PATH}"
    die "${RED}error, failed to change directory to ${repo_name}${NC}"
fi

# commit changes
echo -e "${GREEN}${repo_name}:${NC} git commit -a -m \"new rpc-agents for version ${VERSION}\""
git commit -a -m "new rpc-agents for version ${VERSION}" || {
    cd "${START_PATH}"
    die "${RED}error, failed to commit changes for new rpc-agents version ${VERSION}${NC}"
}

# tag_repo
if [[ -n "${ADD_TAG}" ]]; then
    echo -e "${GREEN}${repo_name}:${NC} git tag -a ${VERSION} -m \"Version ${VERSION}\""
    git tag -a ${VERSION} -m "Version ${VERSION}" || {
        cd "${START_DIR}"
        die "${RED}error, failed to add tag for repo ${repo_name}: git tag -a ${VERSION} -m \"Version ${VERSION}\"${NC}"
    }
fi

# push changes if the destination branch is master only
if [[ "${RC_MASTER_BRANCH_NAME}" == "master" ]]; then
    echo -e "${GREEN}Push repo ${repo_name}${NC}"
    echo -e "${GREEN}${repo_name}:${NC} git push origin ${RC_MASTER_BRANCH_NAME}"
    git push origin ${RC_MASTER_BRANCH_NAME} || {
        cd "${START_DIR}"
        die "${RED}error, failed to: git push origin ${RC_MASTER_BRANCH_NAME}${NC}"
    }

    if [[ -n "${ADD_TAG}" ]]; then
        # push tags
        echo -e "${GREEN}${repo_name}:${NC} git push --tags"
        git push --tags || {
            cd "${START_DIR}"
            die "${RED}error, failed to: \"git push --tags\" in ${repo_name}${NC}"
        }
    fi
fi

cd "${START_DIR}"

repo_name=${REPO_THRIFT_NAME}
cd "${repo_name}"
RC=$?
if [[ ${RC} -ne 0 ]]; then
    cd "${START_PATH}"
    die "${RED}error, failed to change directory to ${repo_name}${NC}"
fi

echo -e "${GREEN}${repo_name}:${NC} git checkout ${RC_MASTER_BRANCH_NAME}"
git checkout "${RC_MASTER_BRANCH_NAME}" || {
    cd "${START_PATH}"
    die "${RED}error, failed to checkout to ${RC_MASTER_BRANCH_NAME} for repo ${repo_name}"
}

echo -e "${GREEN}${repo_name}:${NC} git pull"
git pull || {
    cd "${START_PATH}"
    die "${RED}error, failed to: git pull${NC}"
}

# echo -e "${GREEN}${repo_name}:${NC} git merge --no-ff --no-commit ${RC_BRANCH_NAME}"
# git merge --no-ff --no-commit "${RC_BRANCH_NAME}"
echo -e "${GREEN}${repo_name}:${NC} git merge --no-ff --no-edit ${RC_BRANCH_NAME}"
git merge --no-ff --no-edit "${RC_BRANCH_NAME}" || {
    cd "${START_PATH}"
    die "${RED}error, failed to merge ${RC_BRANCH_NAME} to ${RC_MASTER_BRANCH_NAME}${NC}"
}

# push changes
echo -e "${GREEN}Push repo ${repo_name}${NC}"
echo -e "${GREEN}${repo_name}:${NC} git push origin ${RC_MASTER_BRANCH_NAME}"
git push origin ${RC_MASTER_BRANCH_NAME} || {
    cd "${START_DIR}"
    die "${RED}error, failed to: git push origin ${RC_MASTER_BRANCH_NAME}${NC}"
}

# push tags
if [[ -n "${ADD_TAG}" ]]; then
    echo -e "${GREEN}${repo_name}:${NC} git push --tags"
    git push --tags || {
        cd "${START_DIR}"
        die "${RED}error, failed to: git push --tags in ${repo_name}${NC}"
    }
fi

# go back to initial directory
cd "${START_DIR}"

exit 0
