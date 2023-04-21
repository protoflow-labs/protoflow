package main

import (
	"context"
	"log"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// list namespaces
	list, err := clientset.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, namespace := range list.Items {
		log.Printf("Namespace: %s, Status: %s", namespace.Name, namespace.Status.Phase)
	}

}
