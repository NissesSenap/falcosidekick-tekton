package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

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
	var CriticalNamespaces = []string{"kube-system", "kube-public", "kube-node-lease", "falco"}
	var alert Alert

	bodyReq := os.Getenv("BODY")
	if bodyReq == "" {
		panic("Need to get environment variable BODY")
	}
	bodyReqByte := []byte(bodyReq)
	json.Unmarshal(bodyReqByte, &alert)

	podName := alert.OutputFields.K8SPodName
	namespace := alert.OutputFields.K8SNsName
	log.Printf("PodName: %v & Namespace: %v", podName, namespace)

	log.Printf("Rule: %v", alert.Rule)
	var critical bool
	for _, ns := range CriticalNamespaces {
		if ns == namespace {
			critical = true
			break
		}
	}

	// setup kubeClient
	var kubeClient *kubernetes.Clientset
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	if !critical {
		log.Printf("Deleting pod %s from namespace %s", podName, namespace)
		err := kubeClient.CoreV1().Pods(namespace).Delete(context.Background(), podName, metaV1.DeleteOptions{})
		if err != nil {
			log.Fatalf("Unable to delete pod due to err %v", err)
			os.Exit(1)
		}
	}
}
