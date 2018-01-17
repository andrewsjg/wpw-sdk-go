#!/bin/bash

# colors
typeset RED='\033[0;31m'
typeset GREEN='\033[0;32m'
typeset NC='\033[0m'  # No Color

[[ -z "${MIRRORS}" ]] && export MIRRORS="https://github.com/WPTechInnovation"
#export MIRRORS="file:///c/Users/wabe/clones/release_script/copy_of_repos"
export MIRRORS="file://${HOME}/clones/mirrors"

typeset -r REPO_GO="${MIRRORS}/wpw-sdk-go.git"
typeset -r REPO_DOTNET="${MIRRORS}/wpw-sdk-dotnet.git"
typeset -r REPO_NODEJS="${MIRRORS}/wpw-sdk-nodejs.git"
typeset -r REPO_PYTHON="${MIRRORS}/wpw-sdk-python.git"
typeset -r REPO_JAVA="${MIRRORS}/wpw-sdk-java.git"
typeset -r REPO_IOT="${MIRRORS}/wpw-sdk-iot-core.git"
typeset -r REPO_THRIFT="${MIRRORS}/wpw-sdk-thrift.git"


typeset -r REPO_GO_NAME="wpw-sdk-go"
typeset -r REPO_DOTNET_NAME="wpw-sdk-dotnet"
typeset -r REPO_NODEJS_NAME="wpw-sdk-nodejs"
typeset -r REPO_PYTHON_NAME="wpw-sdk-python"
typeset -r REPO_JAVA_NAME="wpw-sdk-java"
typeset -r REPO_IOT_NAME="wpw-sdk-iot-core"
typeset -r REPO_THRIFT_NAME="wpw-sdk-thrift"
#typeset ALL_REPOS_NAMES="${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} ${REPO_GO_NAME}"

typeset RC_BRANCH_NAME=""
typeset MASTER_BRANCH_NAME=""
typeset VERSION=""
typeset IN_REPOS_NAMES=()
typeset IN_REPOS=()
typeset PUSH=false
typeset PUSH_ONLY=false
typeset CLEAN=false
typeset ADD_TAG=""

# change GOPATH to ./go
export GOPATH=`pwd`/go
export WPW_SDK_GO_PATH=${GOPATH}/src/github.com/WPTechInnovation/${REPO_GO_NAME}

# functions
function cleanup {
    if [[ ${CLEAN} == true ]]; then
        echo -e "${GREEN}****************************************${NC}"
        echo -e "${GREEN}********** Remove directories **********${NC}"
        echo -e "${GREEN}****************************************${NC}"
        echo
        for repo_name in ${ALL_REPOS_NAMES[@]};
        do
            if [ -d "${repo_name}" ]; then
                echo -e "${GREEN} cleanup: Removing directory ${repo_name}${NC}"
                # Control will enter here if $DIRECTORY exists.
                rm -fr "${repo_name}"
            fi
        done
    fi
}

function join_by {
    local IFS="$1"; shift; echo "$*";
}

# input attributes
while true; do
  case "$1" in
    -v | --version ) VERSION="$2"; shift; shift ;;
    -t | --add_tag ) ADD_TAG="-t"; shift ;;
    -b | --branch ) RC_BRANCH_NAME="$2"; shift; shift ;;
    -m | --master_branch ) MASTER_BRANCH_NAME="$2"; shift; shift ;;
    -p | --push ) PUSH=true; shift ;;
    -o | --push_only ) PUSH_ONLY=true; shift ;;
    -c | --clean ) CLEAN=true; shift ;;
    -r | --repos_names )
        IN_REPOS_NAMES=(${2//,/ })
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
        shift ;;
    * ) break ;;
  esac
done

# vfy input attributes
if [[ -z ${VERSION} ]]; then
    echo -e "${RED}error, version name not defined${NC}"
    exit 1
fi

if [[ -z ${RC_BRANCH_NAME} ]]; then
    echo -e "${RED}error, branch for release candidate is not defined, please specify argument for option -b${NC}"
    exit 1
fi

if [[ -z ${MASTER_BRANCH_NAME} ]]; then
    echo -e "${RED}error, master branch not defined, lease specify argument for option -m${NC}"
    exit 1
fi

if [[ ${PUSH} == true && ${PUSH_ONLY} == true ]]; then
    echo -e "${RED}error, both parameters: push (-p) and push_only (-o) cannot be set${NC}"
    exit 1
fi

# determine repos to update
if [[ ${#IN_REPOS_NAMES[@]} -ne 0 ]]; then
    ALL_REPOS_NAMES=("${IN_REPOS_NAMES[@]}")
else
    ALL_REPOS_NAMES=( ${REPO_GO_NAME} ${REPO_DOTNET_NAME} ${REPO_NODEJS_NAME} ${REPO_PYTHON_NAME} ${REPO_JAVA_NAME} ${REPO_IOT_NAME} ${REPO_THRIFT_NAME} )
fi

typeset ALL_REPOS_NAMES_STRING=`join_by , "${ALL_REPOS_NAMES[@]}"`

if [[ ${#IN_REPOS[@]} -ne 0 ]]; then
    ALL_REPOS=("${IN_REPOS[@]}")
else
    ALL_REPOS=( ${REPO_GO} ${REPO_DOTNET} ${REPO_NODEJS} ${REPO_PYTHON} ${REPO_JAVA} ${REPO_IOT} ${REPO_THRIFT} )
fi

typeset ALL_REPOS_STRING=`join_by , "${ALL_REPOS[@]}"`

if [[ ${PUSH_ONLY} == false ]]; then
    # prepare_clones
    echo
    echo -e "${GREEN}****************************************${NC}"
    echo -e "${GREEN}*** Prepare clones (prepare_env.sh). ***${NC}"
    echo -e "${GREEN}****************************************${NC}"
    echo
    ./prepare_clones.sh -b ${RC_BRANCH_NAME} -r ${ALL_REPOS_NAMES_STRING} -e ${ALL_REPOS_STRING}
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to prepares clones${NC}"
        cleanup
        exit 2
    fi

    echo
    echo -e "${GREEN}****************************************${NC}"
    echo -e "${GREEN}*********  prepare submodules  *********${NC}"
    echo -e "${GREEN}****************************************${NC}"
    echo
    ./prepare_submodules.sh -v "${VERSION}" -b "${RC_BRANCH_NAME}" -m "${MASTER_BRANCH_NAME}" "${ADD_TAG}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to prepare submodules${NC}"
        cleanup
        exit 3
    fi

    # update submodules
    echo
    echo -e "${GREEN}************************************************${NC}"
    echo -e "${GREEN}*** Update submodules (update_submodules.sh) ***${NC}"
    echo -e "${GREEN}************************************************${NC}"
    echo
    ./update_submodules.sh -b "${RC_BRANCH_NAME}" -r "${ALL_REPOS_NAMES_STRING}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to update submodules${NC}"
        cleanup
        exit 4
    fi

    # merge release candidate to develop/master
    echo
    echo -e "${GREEN}**********************************************${NC}"
    echo -e "${GREEN}*** Merge release condidate (merge_rc.sh). ***${NC}"
    echo -e "${GREEN}**********************************************${NC}"
    echo
    ./merge_rc.sh -b "${RC_BRANCH_NAME}" -m "${MASTER_BRANCH_NAME}" -r "${ALL_REPOS_NAMES_STRING}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to merge branches${NC}"
        cleanup
        exit 5
    fi

    if [[ -n "${ADD_TAG}" ]]; then
        # tag changes
        echo
        echo -e "${GREEN}***************************************${NC}"
        echo -e "${GREEN}*** Tag repositories (tag_repos.sh) ***${NC}"
        echo -e "${GREEN}***************************************${NC}"
        echo
        ./tag_repos.sh -v "${VERSION}" -r "${ALL_REPOS_NAMES_STRING}"
        RC=$?
        if [[ ${RC} != 0 ]]
        then
            echo -e "${RED}error, failed to tag version${NC}"
            cleanup
            exit 6
        fi
    fi
fi

if [[ ${PUSH} == true || ${PUSH_ONLY} == true ]]; then
    # push
    echo
    echo -e "${GREEN}*****************************************${NC}"
    echo -e "${GREEN}*** Push repositories (push_repos.sh) ***${NC}"
    echo -e "${GREEN}*****************************************${NC}"
    echo
    ./push_repos.sh -m "${MASTER_BRANCH_NAME}" -r "${ALL_REPOS_NAMES_STRING}"
    RC=$?
    if [[ ${RC} != 0 ]]
    then
        echo -e "${RED}error, failed to push changes${NC}"
        cleanup
        exit 7
    fi
fi

cleanup

exit 0
