package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var kubeClient *kubernetes.Clientset

func init() {

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
}

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

var CriticalNamespaces = []string{"kube-system", "kube-public", "kube-node-lease", "falco", "openfaas", "openfaas-fn"}

func main() {
	var alert Alert

	bodyReq := os.Getenv("BODY")
	if bodyReq == "" {
		panic("Need to get ENV var body")
	}
	fmt.Println(bodyReq)
	bodyReqByte := []byte(bodyReq)
	json.Unmarshal(bodyReqByte, alert)

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		json.Unmarshal(body, &alert)

		podName := alert.OutputFields.K8SPodName
		namespace := alert.OutputFields.K8SNsName

		var critical bool
		for _, ns := range CriticalNamespaces {
			if ns == namespace {
				critical = true
				break
			}
		}

		if !critical {
			log.Printf("Deleting pod %s from namespace %s", podName, namespace)
			kubeClient.CoreV1().Pods(namespace).Delete(context.Background(), podName, metaV1.DeleteOptions{})
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
