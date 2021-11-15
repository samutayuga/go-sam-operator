package utils

import (
	"flag"
	"html/template"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	displayTemplate *template.Template
)

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
func InClusterConfig() *rest.Config {
	if config, err := rest.InClusterConfig(); err == nil {
		return config
	} else {
		log.Printf("Error while calling rest.InClusterConfig(), %v\n", err)
	}
	return nil

}
func OutClusterConfig() *rest.Config {
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
func GetIsInClusterEnvVar(inClusterDeployment string) bool {
	if inClusterDeployment != "" {
		if isInCluster, errParse := strconv.ParseBool(inClusterDeployment); errParse == nil {
			return isInCluster
		} else {
			log.Printf("error while converting IS_IN_CLUSTER_DEP, %v\n", errParse)
		}

	}
	return false
}
func InitKubeConnection(isInClusterDep bool) *kubernetes.Clientset {
	log.Printf("this is initialization with IS_IN_CLUSTER_DEP=%v\n", isInClusterDep)
	var kubeConf *rest.Config
	if isInClusterDep {
		kubeConf = InClusterConfig()
	} else {
		kubeConf = OutClusterConfig()
	}
	if kubeConf != nil {
		if clientset, clientSetCreateErr := kubernetes.NewForConfig(kubeConf); clientSetCreateErr == nil {
			return clientset

		} else {
			log.Fatalf("Error while creating the connection %v\n", clientSetCreateErr)
		}

	}
	return nil
}
func ListingFilesInCurrentDir() {
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
}
