package main

import (
	"context"
	"os"
	"testing"

	//"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Test struct {
	createFunc func() error
	deleteFunc func() error
}

var (
	clientset *kubernetes.Clientset
	tests     map[string]Test
)

func init() {
	configPath := clientcmd.RecommendedHomeFile
	if path := os.Getenv(clientcmd.RecommendedConfigPathEnvVar); path != "" {
		configPath = path
	}

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	ns := "default"

	tests = map[string]Test{
		"configmaps": {
			createFunc: func() (err error) {
				_, err = clientset.CoreV1().ConfigMaps(ns).Create(ctx, configMap, metav1.CreateOptions{})
				return
			},
			deleteFunc: func() (err error) {
				err = clientset.CoreV1().ConfigMaps(ns).Delete(ctx, configMap.Name, metav1.DeleteOptions{})
				return
			},
		},
		"secrets": {
			createFunc: func() (err error) {
				_, err = clientset.CoreV1().Secrets(ns).Create(ctx, secret, metav1.CreateOptions{})
				return
			},
			deleteFunc: func() (err error) {
				err = clientset.CoreV1().Secrets(ns).Delete(ctx, secret.Name, metav1.DeleteOptions{})
				return
			},
		},
		"deployments": {
			createFunc: func() (err error) {
				_, err = clientset.AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})
				return
			},
			deleteFunc: func() (err error) {
				err = clientset.AppsV1().Deployments(ns).Delete(ctx, deployment.Name, metav1.DeleteOptions{})
				return
			},
		},
		"services": {
			createFunc: func() (err error) {
				_, err = clientset.CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
				return
			},
			deleteFunc: func() (err error) {
				err = clientset.CoreV1().Services(ns).Delete(ctx, service.Name, metav1.DeleteOptions{})
				return
			},
		},
	}
}

func TestApply(t *testing.T) {
	for name, test := range tests {
		t.Logf("Test %s", name)
		err := test.createFunc()
		if err != nil {
			t.Errorf("Test %s FAILED: %+v", name, err)
		}
	}
}

func TestDelete(t *testing.T) {
	for name, test := range tests {
		t.Logf("Test %s", name)
		err := test.deleteFunc()
		if err != nil {
			t.Errorf("Test %s FAILED: %+v", name, err)
		}
	}
}
