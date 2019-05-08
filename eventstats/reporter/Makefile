.PHONY: all push local_minikube

FULL_IMAGE_NAME=bobais/eventstats-reporter:latest
MINIKUBE_PROFILE?=rhavelka-minikube
LOCAL_REPORTER_IMAGE=bobais/eventstats-watcher:local

all: build

build:
	docker build -t ${FULL_IMAGE_NAME} .

push: build
	docker push ${FULL_IMAGE_NAME}

minikube_build:
	@eval $$(minikube docker-env -p ${MINIKUBE_PROFILE}) ; \
	docker build -t ${LOCAL_REPORTER_IMAGE} .


