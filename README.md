# Cluster EventStats Reporter #

Service consists of 3 components.

* Watcher (Go) - simply watches events with kubernetes go client and generates statistics scrapeable by prometheus
* Reporter (Python) - generates and sends e-mail
* Prometheus

## Run ##

For deployment are prepared helm charts to create Deployment (Watcher), Cron Job (Reporter), Stateful Set (Prometheus) and more.

### Local minikube ###

Prerequisities - ```minikube```, ```docker```, ```helm```, ```make```

Run ```make``` command in the root directory and wait. It will start local
minikube, build containers inside minikube localy, deploy local containers
using helm and opens minikube dashboard. In terminal should be visible link
to prometheus service to check data. Empty REP_SMTP_HOST will cause that
HTML data will be logged - trigger cron job manually and check pod's logs.

**Other ```make``` targets**

* ```deploy_from_registry``` - re-deploys into local minikube using latest images from registry
* ```clean``` - removes minikube

### In cluster with SMTP using helm ###

Prerequisities - ```helm```, SMTP server

```bash
helm install \
    --set reporter.env.REP_SMTP_HOST='smtp.example.org' \
    --set reporter.env.REP_SMTP_HOST_PORT='465' \
    --set reporter.env.REP_SMTP_USER='jon-doe' \
    --set reporter.envsec.REP_SMTP_USER_PASSWORD='password' \
    --set reporter.env.REP_SMTP_SSL='True' \
    --set reporter.env.REP_RECIPIENTS='jon-doe@example.org' \
    --set reporter.env.REP_FROM='jane-doe@example.org' \
    --set reporter.env.REP_TIMEWINDOW='1d' \
    --set reporter.schedule="0 7 * * *" \
    helm/
```

## Reporter notes ##

### Parameters *REP_TIMEWINDOW* and *schedule* ###

*REP_TIMEWINDOW* defines timeframe of interest from NOW()-REP_TIMEWINDOW to NOW(). *schedule* defines the frequency (Cron Job) when reporter will generate and send the report.

## Watcher notes ##

* Written in Go but without proper public repo.
* Dependencies not managed via tool.
* Not 100% accurate.

## Sources ##

* <https://github.com/bobais/k8s-homework>
* <https://hub.docker.com/r/bobais/eventstats-watcher>
* <https://hub.docker.com/r/bobais/eventstats-reporter>
