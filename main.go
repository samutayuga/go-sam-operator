package main

import (
	"context"
	"github.com/samutayuga/go-sam-operator/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
	"os"
	"time"
)

var (
	clientset          *kubernetes.Clientset
	inClusterDep       = false
	monitoredNamespace = "default"
	podInterface       v12.PodInterface
)

func main() {
	log.Println("This is hook first line")

	var timeOut = int64(30)
	if monitoredNamespace != "" {
		log.Printf("Hook is triggered for the namespace %s\n", monitoredNamespace)
		monitor(monitoredNamespace, &timeOut)
	} else {
		log.Println("The hook is invoked without namespace argument")
	}

}
func monitor(namespace string, timeOut *int64) {
	log.Printf("Hook is triggered for the namespace %s and timeout for watching an event %d seconds\n", namespace, *timeOut)
	if clientset != nil {
		podInterface = clientset.CoreV1().Pods(namespace)
		watch, errWatch := podInterface.Watch(context.Background(), metav1.ListOptions{})
		if errWatch != nil {
			log.Printf("No changes of the resources after %d waiting %v\n", *timeOut, errWatch)
			if pods, errList := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{}); errList == nil {
				if len(pods.Items) > 0 {
					utils.DisplayPodList(pods.Items)
				}
			} else {
				log.Printf("Error while pod listing  %v\n", errList)
			}

		} else {
			go func() {
				for anEvent := range watch.ResultChan() {
					if p, ok := anEvent.Object.(*v1.Pod); ok {
						utils.DisplayAPod(p)
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

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	utils.ListingFilesInCurrentDir()
	utils.InitTemplate()

	inClusterDeployment := os.Getenv("IS_IN_CLUSTER_DEP")
	monitoredNamespace = os.Getenv("KUBE_NAMESPACE")
	inClusterDep = utils.GetIsInClusterEnvVar(inClusterDeployment)
	clientset = utils.InitKubeConnection(inClusterDep)
}
