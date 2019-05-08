package utils

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultPort = "10080"

func GetResolvedKubernetesClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Print("Apparently not running in cluster. Will try local $HOME configs.")

		var kubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create in-cluster config: %v", err)
		}
	}

	return kubernetes.NewForConfig(config)
}

func StartHttpServer() {
	log.Printf("Starting prometheus handler: %s", defaultPort)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatalf("Prometheus monitoring failed: %v", http.ListenAndServe(":"+defaultPort, nil))
	}()
}

func GetEventListingOptions() v1.ListOptions {
	timeout := int64(5 * time.Minute)
	listOptions := v1.ListOptions{
		TimeoutSeconds: &timeout,
	}
	return listOptions
}

func PrintIfDebug(debug bool, v ...interface{}) {
	if debug {
		log.Print(v...)
	}
}

func PrintfIfDebug(debug bool, format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}

func IsPositiveStatement(response string) bool {
	switch strings.ToLower(response) {
	case
		"true",
		"yes",
		"1",
		"y":
		return true
	}
	return false
}
