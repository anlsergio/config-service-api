#!/bin/bash

set -e
set -o pipefail

############################
### Script config params ###
############################
# initialize the PID to avoid unintended behavior.
PORT_FORWARDING_PID=0
REGISTRY_HOST=localhost
# Set the timeout to wait for the registry host to
# be available.
REGISTRY_HOST_CHECK_TIMEOUT=5
REGISTRY_PORT=5000

K8S_MANIFESTS_DIR=./k8s
K8S_DEPLOYMENT_NAME=config-app-deployment

LOCAL_APP_IMAGE_NAME=config-service
FULL_IMAGE_NAME=$REGISTRY_HOST:$REGISTRY_PORT/$LOCAL_APP_IMAGE_NAME:latest


# This function blocks the execution until
# the registry host is available for connection.
# If the timeout expires, it will exit the script.
function waitForRegistryToBeAvailable() {
    echo "Waiting for the Registry host to be available..."
    echo "Timeout: $REGISTRY_HOST_CHECK_TIMEOUT seconds"

    while true
    do
      if timeout $REGISTRY_HOST_CHECK_TIMEOUT bash -c "</dev/tcp/$REGISTRY_HOST/$REGISTRY_PORT"
      then
        echo "Registry is available!"
        break
      else
        echo "waiting..."
        sleep 1
      fi
    done
}

# This function creates a background process for
# managing the port forwarding so that the Minikube Registry
# is accessible from the localhost.
function createLocalRegistryPortForwarding() {
      kubectl port-forward --namespace kube-system service/registry $REGISTRY_PORT:80 &
      PORT_FORWARDING_PID=$!

      waitForRegistryToBeAvailable

      if [ $PORT_FORWARDING_PID -ne 0 ]; then
          echo "Port forwarding process $PORT_FORWARDING_PID started."
      else
          echo "[ERROR] Port forwarding process failed to start in the background"
          exit 1
      fi
}

# This function terminates the port-forwarding process
# releasing the port allocated for it.
function stopPortForwarding() {
    if ps -p $PORT_FORWARDING_PID > /dev/null; then
        echo "Stopping port forwarding $PORT_FORWARDING_PID"
        kill $PORT_FORWARDING_PID
    fi
}

# Create the port-forwarding so that the Minikube Registry
# is accessible locally:
if ss -tulnp | grep ":$REGISTRY_PORT" > /dev/null
then
    echo "[ERROR] The port defined for Minikube Registry $REGISTRY_PORT is already in use!"
    exit 1
else
    createLocalRegistryPortForwarding
fi

# Tag the local docker image as the pattern that will be
# uploaded to the registry:
# TODO: make sure the deployment will pull the latest image
# for every new deploy.
docker tag $LOCAL_APP_IMAGE_NAME $FULL_IMAGE_NAME

# Push the image to the Minikube registry
docker push $FULL_IMAGE_NAME

# Delete an existing deployment, if there's any.
if kubectl get deployment $K8S_DEPLOYMENT_NAME
then
  kubectl delete deployment $K8S_DEPLOYMENT_NAME
fi

# Apply the K8s manifests to deploy the application.
kubectl apply -f $K8S_MANIFESTS_DIR

# Kill the background port forwarding process after deploy is concluded.
stopPortForwarding
