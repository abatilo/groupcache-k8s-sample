#!/bin/bash
set -ex

PROJECT_NAME="sample"
REGISTRY_PORT=32121

asdf install
ctlptl get registry $PROJECT_NAME-registry || ctlptl create registry $PROJECT_NAME-registry --port=$REGISTRY_PORT
ctlptl get cluster kind-$PROJECT_NAME || ctlptl create cluster kind --name kind-$PROJECT_NAME --registry=$PROJECT_NAME-registry
tilt up --stream
ctlptl get cluster kind-$PROJECT_NAME && ctlptl delete cluster kind-$PROJECT_NAME
ctlptl get registry $PROJECT_NAME-registry && ctlptl delete registry $PROJECT_NAME-registry
