#!/bin/bash

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function deploy {
    echo "Deploying..."
    #VERSION="pipeline-$CI_PIPELINE_ID"
    VERSION="$BITBUCKET_COMMIT"
    curl -X POST -F token=$INFRA_TOKEN -F "ref=master" -F "variables[SERVICE_NAME]=tasktab" -F "variables[SERVICE_VERSION]=$VERSION" https://gitlab.com/api/v4/projects/6986946/trigger/pipeline
}

function deploy-when-master {
    if [[ "$CI_COMMIT_REF_NAME" == "master" ]]; then
      deploy
    fi
}

function dump-schema {
    docker-compose exec -T db /bin/sh -c "/usr/bin/mysqldump -udev -pdev --no-data dev" | grep -v "Using a password on the command line interface can be insecure" | sed 's/ AUTO_INCREMENT=[0-9]*//g' > migrations/0.sql
}
function build {
    docker build -t tasktab .
}
function up {
    docker-compose up -d
}
function backend {
    DEV_MODE=true go run . www
}
function stop {
    docker-compose stop
}
function default {
    help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}
