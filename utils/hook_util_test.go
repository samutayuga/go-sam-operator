package utils

import (
	assert2 "github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"log"
	"testing"
)

func TestGetIsInClusterEnvVar(t *testing.T) {
	cases := []struct {
		strEnvVar string
		expected  bool
	}{
		{"", false},
		{"true", true},
		{"false", false},
	}
	for _, aCase := range cases {
		got := GetIsInClusterEnvVar(aCase.strEnvVar)
		if got != aCase.expected {
			t.Errorf("GetIsInClusterEnvVar(%q)==%t, expect %t", aCase.strEnvVar, got, aCase.expected)
		}
	}
}

func TestWithAssert(t *testing.T) {
	assert := assert2.New(t)
	assert.Equalf(true, GetIsInClusterEnvVar("true"), "If the IS_IN_CLUSTER_DEP is \"true\" then the isInClusterDeployment is true")
	assert.Equalf(false, GetIsInClusterEnvVar("false"), "If the IS_IN_CLUSTER_DEP is \"false\" then the isInClusterDeployment is false")
	assert.Equalf(false, GetIsInClusterEnvVar(""), "If the IS_IN_CLUSTER_DEP is missing then the isInClusterDeployment is false")
}

func TestWatchPod(t *testing.T) {
	kubeCl := KubeClient{
		Clientset: fake.NewSimpleClientset(&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "influxdb-v2",
				Namespace: "session-1",
				Annotations: map[string]string{
					"ProfefeEnabledAnnotationc": "true",
				},
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
			},
		}),
	}
	kubeCl.watchPod("session-1")

}
func TestListPod(t *testing.T) {
	kubeCl := KubeClient{
		Clientset: fake.NewSimpleClientset(&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "influxdb-v2",
				Namespace: "session-1",
				Annotations: map[string]string{
					"ProfefeEnabledAnnotationc": "true",
				},
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
			},
		}),
	}

	pods := kubeCl.listPods("session-1")
	log.Printf("pod %v", pods)

}
