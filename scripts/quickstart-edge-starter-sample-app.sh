#!/bin/bash
HOME_DIR=$(pwd)

set -e

RUN_QUICKSTART=1
SKIP_PREDIX_SERVICES=false
BUILD_APP=false
function local_read_args() {
  while (( "$#" )); do
  opt="$1"
  case $opt in
    -h|-\?|--\?--help)
      PRINT_USAGE=1
      QUICKSTART_ARGS="$SCRIPT $1"
      break
    ;;
    -b|--branch)
      BRANCH="$2"
      QUICKSTART_ARGS+=" $1 $2"
      shift
    ;;
    --skip-predix-services)
      SKIP_PREDIX_SERVICES=true
    ;;
    --build-app)
      BUILD_APP=true
    ;;
    -o|--override)
      QUICKSTART_ARGS=" $SCRIPT"
    ;;
    --skip-setup)
      SKIP_SETUP=true
    ;;
    *)
      QUICKSTART_ARGS+=" $1"
      #echo $1
    ;;
  esac
  shift
  done

  if [[ -z $BRANCH ]]; then
    echo "Usage: $0 -b/--branch <branch> [--skip-setup]"
    exit 1
  fi
}

# default settings
BRANCH="master"
PRINT_USAGE=0
SKIP_SETUP=false
VERSION_JSON="version.json"
PREDIX_SCRIPTS_ORG="PredixDev"
PREDIX_SCRIPTS="predix-scripts"
GITHUB_RAW="https://raw.githubusercontent.com"
IZON_SH="https://raw.githubusercontent.com/PredixDev/izon/1.5.0/izon2.sh"

GITHUB_ORG="PredixDev"
REPO_NAME="predix-edge-sample-scaler-go"
DOCKER_STACK_NAME="predix-edge-sample-scaler-go"
SCRIPT="-script edge-starter-deploy.sh -script-readargs edge-starter-deploy-readargs.sh"
SCRIPT_NAME="quickstart-edge-starter-sample-app.sh"
APP_DIR="edge-sample-scaler-go"
APP_NAME="Predix Edge Sample Scaler App - Go"
TOOLS="Cloud Foundry CLI, Docker, Git, jq, yq, Go, Predix CLI"
TOOLS_SWITCHES="--cf --docker --git --jq --yq --go --predixcli"

TIMESERIES_CHART_ONLY="true"

# Process switches
local_read_args $@

#variables after processing switches
SCRIPT_LOC="$GITHUB_RAW/$GITHUB_ORG/$REPO_NAME/$BRANCH/scripts/$SCRIPT_NAME"
VERSION_JSON_URL="$GITHUB_RAW/$GITHUB_ORG/$REPO_NAME/$BRANCH/version.json"

if [[ "$SKIP_PREDIX_SERVICES" == "true" ]]; then
  QUICKSTART_ARGS="$QUICKSTART_ARGS $SCRIPT -uaa -ts -psts --run-edge-app -repo-name $REPO_NAME -app-name $DOCKER_STACK_NAME"
else
  PS_MODEL="sample-data/predix-asset/Compressor-CMMS-Compressor-2018-Scaled.json"
  QUICKSTART_ARGS="$QUICKSTART_ARGS $SCRIPT -uaa -ts -psts -psmodel $PS_MODEL --run-edge-app -repo-name $REPO_NAME -app-name $DOCKER_STACK_NAME"
fi

function check_internet() {
  set +e
  echo ""
  echo "Checking internet connection..."
  curl "http://google.com" > /dev/null 2>&1
  if [ $? -ne 0 ]; then
    echo "Unable to connect to internet, make sure you are connected to a network and check your proxy settings if behind a corporate proxy"
    echo "If you are behind a corporate proxy, set the 'http_proxy' and 'https_proxy' environment variables.   Please read this tutorial for detailed info about setting your proxy https://www.predix.io/resources/tutorials/tutorial-details.html?tutorial_id=1565"
    exit 1
  fi
  echo "OK"
  echo ""
  set -e
}

function init() {
  currentDir=$(pwd)
  if [[ $currentDir == *"scripts" ]]; then
    echo 'Please launch the script from the root dir of the project'
    exit 1
  fi

  check_internet

  #get the script that reads version.json
  eval "$(curl -s -L $IZON_SH)"

  #download script and cd
  getUsingCurl $SCRIPT_LOC
  chmod 755 $SCRIPT_NAME;
  if [[ ! $currentDir == *"$REPO_NAME" ]]; then
    mkdir -p $APP_DIR
    cd $APP_DIR
  fi

  getVersionFile
  getLocalSetupFuncs $GITHUB_RAW $PREDIX_SCRIPTS_ORG
}

if [[ $PRINT_USAGE == 1 ]]; then
  init
  __print_out_standard_usage
else
  if $SKIP_SETUP; then
    init
  else
    init
    __standard_mac_initialization
  fi
fi

getPredixScripts
#clone the repo itself if running from oneclick script
getCurrentRepo

if ! [ -d "$logDir" ]; then
  mkdir "$logDir"
  chmod 744 "$logDir"
fi

########### custom logic starts here ###########
if [[ "$BUILD_APP" == "true" ]]; then
  cd $PREDIX_SCRIPTS/$REPO_NAME
  dockerVersion=$(grep version Dockerfile | awk -F"=" '{print $2}' | tr -d "\"")
  echo "$dockerVersion"
  docker build --build-arg HTTP_PROXY --build-arg HTTPS_PROXY --build-arg http_proxy --build-arg https_proxy -t predixedge/predix-edge-sample-scaler-go:latest .
  cd ../..
fi

docker service ls
echo ""
echo ""
docker network ls
########### custom logic ends here ###########

echo "quickstart_args=$QUICKSTART_ARGS"
source $PREDIX_SCRIPTS/bash/quickstart.sh $QUICKSTART_ARGS

cat $SUMMARY_TEXTFILE
if [[ $SKIP_PREDIX_SERVICES == false ]]; then
  __append_new_line_log "To see the data in the cloud, using a browser, open the Front-end App URL shown above.  With login=app_user_1, password=App_User_111" "$quickstartLogDir"
fi
__append_new_line_log "" "$logDir"
__append_new_line_log "Successfully completed $APP_NAME installation!" "$quickstartLogDir"
__append_new_line_log "" "$logDir"
