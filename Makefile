.PHONY: all deploy_from_registry open_dashboard print_services clean

MINIKUBE_PROFILE?=rhavelka-minikube
LOCAL_REPORTER_IMAGE=bobais/eventstats-reporter:local
LOCAL_WATCHER_IMAGE=bobais/eventstats-watcher:local
LOCAL_HELM_APP_NAME=rhavelka
REPORTER_IMAGE=bobais/eventstats-reporter:latest
WATCHER_IMAGE=bobais/eventstats-watcher:latest

EXECUTABLES = minikube docker helm
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))


all: deploy_images print_services open_dashboard

start_minikube:
	@minikube status -p ${MINIKUBE_PROFILE} || \
		minikube start -p ${MINIKUBE_PROFILE} --cpus 4 --memory 4096

init_helm: start_minikube
	@helm version || \
		( helm init && \
		  until helm version 2>/dev/null; do echo "Waiting for tiller"; sleep 3; done;)

build_images: start_minikube build_watcher build_reporter

build_reporter:
	@eval $$(minikube docker-env -p ${MINIKUBE_PROFILE}) ; \
	docker build -t ${LOCAL_REPORTER_IMAGE} ./eventstats/reporter/

build_watcher:
	@eval $$(minikube docker-env -p ${MINIKUBE_PROFILE}) ; \
	docker build -t ${LOCAL_WATCHER_IMAGE} ./eventstats/watcher/

deploy_images: init_helm purge_deployment build_images
	helm install \
		--name ${LOCAL_HELM_APP_NAME} \
		--set watcher.imagePullPolicy=Never \
		--set watcher.image=${LOCAL_WATCHER_IMAGE} \
		--set reporter.imagePullPolicy=Never \
		--set reporter.image=${LOCAL_REPORTER_IMAGE} \
		--set prometheus.nodePort=31111 \
		--set reporter.env.REP_SMTP_HOST= \
	 helm/

deploy_from_registry: init_helm purge_deployment
	helm install \
		--name ${LOCAL_HELM_APP_NAME} \
		--set watcher.debug=true \
		--set watcher.imagePullPolicy=Always \
		--set watcher.image=${WATCHER_IMAGE} \
		--set reporter.imagePullPolicy=Always \
		--set reporter.image=${REPORTER_IMAGE} \
		--set prometheus.nodePort=31111 \
		--set reporter.env.REP_SMTP_HOST= \
	 helm/

purge_deployment: init_helm
	helm list | grep ${LOCAL_HELM_APP_NAME} && \
		helm delete --purge ${LOCAL_HELM_APP_NAME} && sleep 15 || echo "Deployment not found"

print_services: start_minikube
	@minikube service list -p ${MINIKUBE_PROFILE} | grep ${LOCAL_HELM_APP_NAME}

open_dashboard:
	@minikube dashboard -p ${MINIKUBE_PROFILE}

clean:
	@minikube delete -p ${MINIKUBE_PROFILE}
