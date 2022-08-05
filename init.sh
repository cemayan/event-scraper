#!/bin/bash

SCRIPT_DIR="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
function log_blue() { printf "\x1B[94m>> $1\x1B[39m\n"; }

function install_postgresql() {
   log_blue "postgresql installing..."
   kubectl apply -f https://raw.githubusercontent.com/reactive-tech/kubegres/v1.15/kubegres.yaml
   kubectl apply -f k8s/postgres-secret.yaml
   kubectl apply -f k8s/postgres-config.yaml
   kubectl apply -f k8s/postgres.yaml
}


function install_rabbitmq() {
  log_blue "rabbitmq installing..."
  kubectl krew install rabbitmq
  kubectl rabbitmq create esmq
}
