.PHONY: all push run_against_local_minikube

FULL_IMAGE_NAME="bobais/eventstats-watcher:latest"
APP_PORT=13009
MINIKUBE_PROFILE?=rhavelka-minikube
LOCAL_WATCHER_IMAGE=bobais/eventstats-watcher:local

all: build

build:
	docker build -t ${FULL_IMAGE_NAME} .

push: build
	docker push ${FULL_IMAGE_NAME}

minikube_build:
	@eval $$(minikube docker-env -p ${MINIKUBE_PROFILE}) ; \
	docker build -t ${LOCAL_WATCHER_IMAGE} .

run_against_local_minikube: build
	@echo "Open browser http://localhost:${APP_PORT}/metrics"
	docker run -v ~/.kube:/root/.kube -v ~/.minikube:${HOME}/.minikube -p${APP_PORT}:10080 ${FULL_IMAGE_NAME}
