/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	//"flag"
	"fmt"
	//"k8s.io/client-go/tools/clientcmd"
	"os"
	//"path/filepath"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"errors"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func getConfigmapValue(key string) (error,string) {
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
	namespace := getEnv("DEVOPS_ONDUTY_NAMESPACE","default")
	configmap := getEnv("DEVOPS_ONDUTY_CONFIGMAP","devops-shifts-board")
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configmap, metav1.GetOptions{})
	if kerrors.IsNotFound(err) {
		errMsg := fmt.Sprintf("Configmap %s in namespace %s not found\n", configmap, namespace)
		fmt.Printf(errMsg)
		return errors.New(errMsg),""
	} else if statusError, isStatus := err.(*kerrors.StatusError); isStatus {
		errMsg := fmt.Sprintf("Error getting configmap %s in namespace %s: %v\n",
			configmap, namespace, statusError.ErrStatus.Message)
		fmt.Printf(errMsg)
		return errors.New(errMsg),""
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found configmap %s in namespace %s\n", configmap, namespace)
		value := cm.Data[key]
		if value != "" {
			fmt.Printf("Found value %s for key %s\n", value, key)
			return nil,value
		}
		fmt.Printf("requested key %s not found in configmap %s\n",key,configmap)
		return errors.New(fmt.Sprintf("requested key %s not found in configmap %s",key,configmap)),""
	}
}

func checkForSpecialChange(date time.Time) (error,string) {
	value := date.Format("02-01-2006")
	fmt.Printf("checkForSpecialChange on %s\n", value)
	return getConfigmapValue(value)
}

func checkDefaultSchedule(date time.Time) (error,string) {
	value := date.Weekday().String()
	fmt.Printf("checkDefaultSchedule on %s\n", value)
	return getConfigmapValue(value)
}

func whoIsOnDutyNow() (error,string) {
	now := time.Now()
	var onDuty string
	err,onDuty := checkForSpecialChange(now)
	if err != nil {
		err,onDuty = checkDefaultSchedule(now)
		if err != nil {
			return err,""
		}
	}
	return nil,onDuty
}