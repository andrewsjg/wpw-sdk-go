#!/bin/bash

# colors
typeset RED='\033[0;31m'
typeset GREEN='\033[0;32m'
typeset NC='\033[0m'  # No Color

typeset -r REPO_DOTNET_NAME="wpw-sdk-dotnet"
typeset -r REPO_NODEJS_NAME="wpw-sdk-nodejs"
typeset -r REPO_PYTHON_NAME="wpw-sdk-python"
typeset -r REPO_JAVA_NAME="wpw-sdk-java"
#typeset -r REPO_IOT_NAME="wpw-sdk-iot-core"
#typeset -r REPO_THRIFT_NAME="wpw-sdk-thrift"
typeset -r REPO_GO_NAME="wpw-sdk-go"
#typeset ALL_REPOS_NAMES="${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} ${REPO_GO_NAME}"

typeset RC_BRANCH_NAME=""
typeset RC_MASTER_BRANCH_NAME=""


# git checkout test_branch
# git pull 
# git checkout master
# git pull
# git merge --no-ff --no-commit test_branch
## check for confilcts 
# git commit -m 'merge test_branch branch'
# git push



function cleanup {
    echo
}

function die {
    echo -e "${1}" >&2
    cleanup
    exit 1
}

function usage {
    echo
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -b | --branch )
        RC_BRANCH_NAME="$2"
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
        usage
        exit 1
        ;;
  esac
  shift
done

if [[ -z ${RC_BRANCH_NAME} ]]; then
    die "${RED}error, branch name not defined${NC}"
fi

if [[ -z ${RC_MASTER_BRANCH_NAME} ]]; then
    die "${RED}error, master branch name not defined${NC}"
fi

if [[ ${#IN_REPOS_NAMES[@]} -ne 0 ]]; then
    ALL_REPOS_NAMES=("${IN_REPOS_NAMES[@]}")
else
    ALL_REPOS_NAMES=( ${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_GO_NAME} )
fi

echo -e "${GREEN}Tag repos with name: ${VERSION}.${NC}"
# commit repos
for repo_name in ${ALL_REPOS_NAMES[@]};
do
    CURRENT_PATH=`pwd`
    case "${repo_name}" in
        ${REPO_GO_NAME} )
            cd "${WPW_SDK_GO_PATH}/applications/rpc-agent"
            ;;
        ${REPO_DOTNET_NAME} | ${REPO_NODEJS_NAME} | ${REPO_PYTHON_NAME} | ${REPO_JAVA_NAME} )
            cd "${repo_name}"
            ;;
        * )
            continue
            ;;
    esac

    RC=$?
    if [[ ${RC} -ne 0 ]]; then
        cd "${CURRENT_PATH}"
        die "${RED}error, failed to change directory to ${repo_name}${NC}"
    fi

    # 1. git checkout test_branch (should be already done)
    # 2. git pull (it's not required, just cloned)
    # 3. git checkout master
    # 4. git pull
    # 5. git merge --no-ff --no-edit test_branch

    echo -e "${GREEN}${repo_name}:${NC} git checkout ${RC_MASTER_BRANCH_NAME}"
    git checkout "${RC_MASTER_BRANCH_NAME}"
    RC=$?
    if [[ ${RC} -ne 0 ]]; then
        cd "${CURRENT_PATH}"
        die "${RED}error, failed to checkout to ${RC_MASTER_BRANCH_NAME} for repo ${repo_name}"
    fi

    echo -e "${GREEN}${repo_name}:${NC} git pull"
    git pull
    RC=$?
    if [[ ${RC} -ne 0 ]]; then
        cd "${CURRENT_PATH}"
        die "${RED}error, failed to: git pull${NC}"
    fi

    # echo -e "${GREEN}${repo_name}:${NC} git merge --no-ff --no-commit ${RC_BRANCH_NAME}"
    # git merge --no-ff --no-commit "${RC_BRANCH_NAME}"
    echo -e "${GREEN}${repo_name}:${NC} git merge --no-ff --no-edit ${RC_BRANCH_NAME}"
    git merge --no-ff --no-edit "${RC_BRANCH_NAME}"
    RC=$?
    if [[ ${RC} -ne 0 ]]; then
        cd "${CURRENT_PATH}"
        die "${RED}error, failed to merge ${RC_BRANCH_NAME} to ${RC_MASTER_BRANCH_NAME}${NC}"
    fi

    cd "${CURRENT_PATH}"
done

exit 0
