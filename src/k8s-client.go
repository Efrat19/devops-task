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
	"bytes"
	"context"
	"fmt"
	"github.com/prometheus/common/log"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"time"
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

func getTimeSince(startTime time.Time) time.Duration {
	return time.Now().Sub(startTime)
}

type PodInfo struct{
	Name    string
	Uptime	time.Duration
	Version string
}

var (
	namespace = "colors-2"
	versionPath = "/api/version"
)

func getVersionOf(service *v1.Service) (string,error) {
	if len(service.Spec.Ports) == 0 {
		return "no version",nil
	}
	port := service.Spec.Ports[0].Port
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d%s",service.Name,service.Namespace,port,versionPath)
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Info(bodyString)
		return bodyString,nil
	} else {
		return "",fmt.Errorf("service %s %s request returned %d code",service.Name,versionPath,resp.StatusCode)
	}
}


func getPodsOf(service *v1.Service,clientset *kubernetes.Clientset) (*[]v1.Pod,error) {
	labelSelector := labels.NewSelector()
	for labelKey, labelVal := range service.Spec.Selector {
		requirement, err := labels.NewRequirement(labelKey, selection.Equals, []string{labelVal})
		if err != nil {
			log.Error(err, "Unable to create selector requirement")
		}
		labelSelector = labelSelector.Add(*requirement)
	}
	runningnlyRule := "status.phase=Running"
	options := metav1.ListOptions{
		LabelSelector:       labelSelector.String(),
		FieldSelector:       runningnlyRule,
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), options)
	return &pods.Items,err
}

func getClientSet() (*kubernetes.Clientset,error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return kubernetes.NewForConfig(config)
}

func getPodInfoList() (*[]PodInfo,error) {
	clientset,err := getClientSet()
	var podInfoList = []PodInfo{}
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _,service := range services.Items {
		version,err := getVersionOf(&service)
		if err != nil {
			panic(err.Error())
		}
		servicePods,err := getPodsOf(&service,clientset)
		if err != nil {
			panic(err.Error())
		}
		for _,pod := range *servicePods {
			podInfoList = append(podInfoList, PodInfo{
				Name:    pod.Name,
				Uptime:  getTimeSince(pod.Status.StartTime.Time),
				Version: version,
			})
		}
	}
	return &podInfoList,err
}


func getServiceLog(tail int64,serviceName string) string {
	clientset,err := getClientSet()
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	servicePods,err := getPodsOf(service,clientset)
	if err != nil {
		panic(err.Error())
	}
	if len(*servicePods) == 0 {
		return "no pods"
	}
	podLogOpts := v1.PodLogOptions{TailLines:&tail}
	req := clientset.CoreV1().Pods((*servicePods)[0].Namespace).GetLogs((*servicePods)[0].Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}
