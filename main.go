package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	clientset       *kubernetes.Clientset
	displayTemplate *template.Template
)

func main() {
	fmt.Println("hello")

	namespace := os.Getenv("KUBE_NAMESPACE")
	var timeOut = int64(30)
	if namespace != "" {
		log.Printf("Hook is triggered for the namespace %s\n", namespace)
		monitor(namespace, &timeOut)
	} else {
		log.Println("The hook.sh is invoked without namespace argument")
	}

}
func monitor(namespace string, timeOut *int64) {
	log.Printf("Hook is triggered for the namespace %s and timeout for watching an event %d seconds\n", namespace, *timeOut)
	if clientset != nil {
		watch, errWatch := clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{Watch: true, TimeoutSeconds: timeOut})
		if errWatch != nil {
			log.Printf("No changes of the resources after %d waiting %v\n", timeOut, errWatch)
			if pods, errList := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{}); errList == nil {
				if len(pods.Items) > 0 {
					DisplayPodList(pods.Items)
				}
			}

		} else {
			go func() {
				for anEvent := range watch.ResultChan() {
					if p, ok := anEvent.Object.(*v1.Pod); ok {
						DisplayAPod(p)
					} else {
						log.Fatal("unexpected type")
					}
				}
			}()
		}
		time.Sleep(5 * time.Second)
	} else {
		log.Printf("Cannot call kubernetes api because client set is null")
	}

}

//InitTemplate ...
func InitTemplate() {
	var err error
	displayTemplate, err = template.ParseGlob("*.gotmpl")
	//remember to remove
	if err != nil {
		log.Fatalf("error while parsing the files %s %v", "*.gotmpl", err)
	}
}
func DisplayPodList(pods []v1.Pod) {
	if err := displayTemplate.ExecuteTemplate(os.Stdout, "sample.gotmpl", pods); err != nil {
		log.Printf("error while executing template %s for data %v %v", "sample.gotmpl", pods, err)
	}
}
func DisplayAPod(pod *v1.Pod) {
	if err := displayTemplate.ExecuteTemplate(os.Stdout, "pod_displayer.gotmpl", pod); err != nil {
		log.Printf("error while executing template %s for data %v %v", "pod_displayer.gotmpl", pod, err)
	}
}
func inClusterConfig() *rest.Config {
	if config, err := rest.InClusterConfig(); err == nil {
		return config
	} else {
		log.Printf("Error while calling rest.InClusterConfig(), %v\n", err)
	}
	return nil

}
func outClusterConfig() *rest.Config {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

	}
	flag.Parse()

	if config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {

		log.Printf("Error while calling clientcmd.BuildConfigFromFlags, %v\n", err)
		return nil
	} else {
		return config
	}
}
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if currFolder, errGetFolder := os.Getwd(); errGetFolder == nil {
		log.Printf("Current folder %v\n", currFolder)
		if files, errReadDir := ioutil.ReadDir(currFolder); errReadDir == nil {
			for _, aFile := range files {
				log.Printf("File %v\n", aFile.Name())
			}
		} else {
			log.Panicf("error while getting listing current folder %v\n", errReadDir)
		}
	} else {
		log.Panicf("error while getting current folder %v\n", errGetFolder)
	}
	InitTemplate()

	inClusterDeployment := os.Getenv("IS_IN_CLUSTER_DEP")
	log.Printf("this is initialization with IS_IN_CLUSTER_DEP=%v\n", inClusterDeployment)
	var kubeConf *rest.Config
	if inClusterDeployment == "" {
		kubeConf = outClusterConfig()
	} else {
		if isInCluster, errParse := strconv.ParseBool(inClusterDeployment); errParse == nil {
			if isInCluster {
				kubeConf = inClusterConfig()
			} else {
				kubeConf = outClusterConfig()
			}

		} else {
			log.Printf("error while converting IS_IN_CLUSTER_DEP, %v\n", errParse)
		}
	}
	if kubeConf != nil {
		var clientSetCreateErr error
		if clientset, clientSetCreateErr = kubernetes.NewForConfig(kubeConf); clientSetCreateErr != nil {
			log.Fatalf("Error while creating the connection %v\n", clientSetCreateErr)
		}
	}

}
