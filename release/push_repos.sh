#!/bin/bash

# colors
typeset RED='\033[0;31m'
typeset GREEN='\033[0;32m'
typeset NC='\033[0m'  # No Color

typeset -r REPO_GO_NAME="wpw-sdk-go"
typeset -r REPO_DOTNET_NAME="wpw-sdk-dotnet"
typeset -r REPO_NODEJS_NAME="wpw-sdk-nodejs"
typeset -r REPO_PYTHON_NAME="wpw-sdk-python"
typeset -r REPO_JAVA_NAME="wpw-sdk-java"
typeset -r REPO_IOT_NAME="wpw-sdk-iot-core"
typeset -r REPO_THRIFT_NAME="wpw-sdk-thrift"

#typeset ALL_REPOS_NAMES="${REPO_GO_NAME} ${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME}"

typeset RC_MASTER_BRANCH_NAME=""

function cleanup {
    echo
}

function die {
    echo -e "${1}" >&2
    cleanup
    exit 1
}

function usage {
    echo "usage: push_repos.sh -m|--master_branch <destination branch>"
    echo "                     [-r|--repos_names <coma separated repos name>]"
    echo "                     [-n|--no-color]"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
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

if [[ -z ${RC_MASTER_BRANCH_NAME} ]]; then
    die "${RED}error, master branch name not defined${NC}"
fi

if [[ ${#IN_REPOS_NAMES[@]} -ne 0 ]]; then
    ALL_REPOS_NAMES=("${IN_REPOS_NAMES[@]}")
else
    ALL_REPOS_NAMES=( ${REPO_GO_NAME} ${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} )
fi

echo -e "${GREEN}Push repos.${NC}"
START_DIR=`pwd`

# commit repos
for repo_name in ${ALL_REPOS_NAMES[@]};
do
    case "${repo_name}" in
        ${REPO_GO_NAME} )
            cd ${WPW_SDK_GO_PATH}
            RC=$?
            ;;
        ${REPO_DOTNET_NAME} |  ${REPO_NODEJS_NAME} | ${REPO_PYTHON_NAME} | ${REPO_JAVA_NAME} )
            cd ${repo_name}
            RC=$?
            ;;
        * )
            continue
            ;;
    esac

    if [[ ${RC} -ne 0 ]]; then
        cd "${START_DIR}"
        die "${RED}error, failed to change directory to ${repo_name}${NC}"
    fi

    # vfy if branch name is correct
    CURRENT_BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
    if [[ ${CURRENT_BRANCH_NAME} != "${RC_MASTER_BRANCH_NAME}" ]]; then
        cd "${START_DIR}"
        die "${RED}error, current branch name ${CURRENT_BRANCH_NAME} is different than ${RC_MASTER_BRANCH_NAME} for ${repo_name}${NC}"
    fi

    echo -e "${GREEN}${repo_name}:${NC} git push origin ${RC_MASTER_BRANCH_NAME}"
    git push origin ${RC_MASTER_BRANCH_NAME} || {
        cd "${START_DIR}"
        die "${RED}error, failed to: git push origin ${RC_MASTER_BRANCH_NAME}${NC}"
    }

    echo -e "${GREEN}${repo_name}:${NC} git push --tags"
    git push --tags || {
        cd "${START_DIR}"
        die "${RED}error, failed to: git push --tags in ${repo_name}${NC}"
    }

    cd "${START_DIR}"
done

exit 0
