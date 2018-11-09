package main

import "C"
import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

//export Init
func Init(namespace string, confimap string) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()

	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	watcher, err := clientset.CoreV1().ConfigMaps(namespace).Watch(v1.ListOptions{})

	for {
		if err != nil {
			fmt.Printf("Exception on GO Kubernetes Client: %s", err)
		}

		<-watcher.ResultChan()

		config, err := clientset.CoreV1().ConfigMaps(namespace).Get(confimap, v1.GetOptions{})

		file, errort := os.Open("server.cfg")

		if err != nil && errort != nil {
			servers := config.Data["server.cfg"]

			file.WriteString(servers)

			file.Close()
		}
	}
}

//export Close
func Close() {

}
func main() {

}
