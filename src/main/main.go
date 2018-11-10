package main

import (
	"encoding/base64"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"strconv"
	"unicode/utf8"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func Init(namespace string, configMap string) {
	log.Printf("Connecting to Kubernetes.")
	debug, berr := strconv.ParseBool(os.Getenv("DEBUG"))

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

	watcher, err := clientset.CoreV1().ConfigMaps(namespace).Watch(v1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", configMap).String(),
		Watch:         true,
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Starting watch Config Map: %s", configMap)

	ch := watcher.ResultChan()

	for event := range ch {
		config, ok := event.Object.(*v12.ConfigMap)

		if !ok {
			log.Printf("Exception on GO Kubernetes Client: can't get ConfigMap")
			continue
		}

		file, errort := os.OpenFile("server.cfg", os.O_WRONLY, 0660)

		if os.IsNotExist(errort) {
			log.Printf("Creating new server.cfg file on system")
			file, errort = os.Create("server.cfg")

			if errort != nil {
				log.Printf("Can't create file on system %s", err)
			}
		}

		servers := config.Data["server.cfg"]

		if len(servers) == 0 {
			log.Printf("New server.cfg is empty!")
		}

		if errort == nil {
			if berr == nil && debug {
				log.Printf("Change detected on %s, handling %s", configMap, event.Type)

				log.Printf("new config: %s", servers)
			}

			decoded, erro := base64.StdEncoding.DecodeString(servers)

			if erro == nil {
				if _, erro := file.Write(decoded); erro != nil {
					log.Printf("Can't write to file: %s", erro.Error())
				}

				if erro := file.Close(); erro != nil {
					log.Printf("Can't close file: %s", erro.Error())
				}
			} else {
				log.Printf("Error while decoding the config: %s", erro.Error())
			}
		} else {
			log.Printf("Error on write to file: %s", errort.Error())
		}
	}
}

func main() {
	log.Printf("Initialing Go application.")

	configmap := os.Getenv("CONFIGMAP")
	namespace := os.Getenv("NAMESPACE")

	if utf8.RuneCountInString(configmap) > 0 {
		Init(namespace, configmap)
	}
}
