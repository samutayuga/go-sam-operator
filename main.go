package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	clientset       *kubernetes.Clientset
	displayTemplate *template.Template
)

func main() {
	fmt.Println("hello")
	namespace := flag.Arg(0)
	var timeOut = int64(30)
	if namespace != "" {
		log.Printf("Hook is triggered for the namespace %s\n", namespace)
		monitor(namespace, &timeOut)
	} else {
		log.Println("The hook is invoked without namespace argument")
	}

}
func monitor(namespace string, timeOut *int64) {
	log.Printf("Hook is triggered for the namespace %s and timeout for watching an event %d seconds\n", namespace, *timeOut)
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
func init() {
	InitTemplate()
	fmt.Println("this is initialization")
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)

}
