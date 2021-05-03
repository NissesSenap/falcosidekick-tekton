package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Alert falco data structure
type Alert struct {
	Output       string    `json:"output"`
	Priority     string    `json:"priority"`
	Rule         string    `json:"rule"`
	Time         time.Time `json:"time"`
	OutputFields struct {
		ContainerID              string      `json:"container.id"`
		ContainerImageRepository interface{} `json:"container.image.repository"`
		ContainerImageTag        interface{} `json:"container.image.tag"`
		EvtTime                  int64       `json:"evt.time"`
		FdName                   string      `json:"fd.name"`
		K8SNsName                string      `json:"k8s.ns.name"`
		K8SPodName               string      `json:"k8s.pod.name"`
		ProcCmdline              string      `json:"proc.cmdline"`
	} `json:"output_fields"`
}

func main() {
	criticalNamespaces := map[string]bool{
		"kube-system":     true,
		"kube-public":     true,
		"kube-node-lease": true,
		"falco":           true,
	}
	var falcoEvent Alert

	bodyReq := os.Getenv("BODY")
	if bodyReq == "" {
		panic("Need to get environment variable BODY")
	}
	bodyReqByte := []byte(bodyReq)
	err := json.Unmarshal(bodyReqByte, &falcoEvent)
	if err != nil {
		panic(fmt.Errorf("The data doesent match the struct %w", err))
	}

	kubeClient, err := setupK8sClient()
	if err != nil {
		panic(fmt.Errorf("Unable to create in-cluster config: %w", err))
	}

	err = deletePod(kubeClient, falcoEvent, criticalNamespaces)
	if err != nil {
		log.Fatalf("Unable to delete pod due to err %v", err)
	}
}

// setupK8sClient
func setupK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

// deletePod, if not part of the criticalNamespaces the pod will be deleted
func deletePod(kubeClient *kubernetes.Clientset, falcoEvent Alert, criticalNamespaces map[string]bool) error {
	podName := falcoEvent.OutputFields.K8SPodName
	namespace := falcoEvent.OutputFields.K8SNsName
	log.Printf("PodName: %v & Namespace: %v", podName, namespace)

	log.Printf("Rule: %v", falcoEvent.Rule)
	for ns := range criticalNamespaces {
		if criticalNamespaces[ns] {
			log.Printf("The pod %v won't be deleted due to it's part of the critical ns list: %v ", podName, ns)
			return nil
		}
	}

	log.Printf("Deleting pod %s from namespace %s", podName, namespace)
	err := kubeClient.CoreV1().Pods(namespace).Delete(context.Background(), podName, metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
