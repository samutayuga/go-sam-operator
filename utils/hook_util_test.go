package utils

import (
	"context"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"testing"
)

type KubeApiMock struct {
	mock.Mock
}

func (k *KubeApiMock) Watch(context context.Context, options v1.ListOptions) (watch.Interface, error) {
	args := k.Called(context, options)
	return args.Get(0).(watch.Interface), args.Error(1)
}
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
