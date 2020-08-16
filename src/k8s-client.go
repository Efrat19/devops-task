package main

import (
	"bytes"
	"context"
	"fmt"
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
	"github.com/apex/log"
)


type PodInfo struct{
	Name    string
	Uptime	time.Duration
	Version string
}

const (
	namespace = "k-bot-test"
	versionPath = "/api/version"
)


func getPodInfoList() (*[]PodInfo,error) {
	log.Debugf("getPodInfoList called")
	clientset,err := getClientSet()
	if err != nil {
		log.Error("Unable to getClientSet")
		return nil, err
	}
	var podInfoList = []PodInfo{}
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error("Unable to list services")
		return nil, err
	}
	for _,service := range services.Items {
		version,err := getServiceVersion(&service)
		if err != nil {
			log.Errorf("Unable to get service version\n%v",err)
			continue
		}
		servicePods,err := getServicePods(&service,clientset)
		if err != nil {
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


func getServiceLog(tail int64,serviceName string) (string,error) {
	log.Debugf("getServiceLog called with tail %d and serviceName %s",tail,serviceName)
	clientset,err := getClientSet()
	if err != nil {
		log.Error("Unable to getClientSet")
		return "", err
	}
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Error("Unable to list services")
		return "", err
	}
	servicePods,err := getServicePods(service,clientset)
	if err != nil {
		log.Error("Unable to get service pods")
		return "", err
	}
	if len(*servicePods) == 0 {
		log.Error("service has no pods")
		return "", fmt.Errorf("no pods found for service %s",serviceName)
	}
	podLogOpts := v1.PodLogOptions{TailLines:&tail}
	req := clientset.CoreV1().Pods((*servicePods)[0].Namespace).GetLogs((*servicePods)[0].Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		log.Error("error in opening stream")
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		log.Error("error in copy information from podLogs to buf")
		return "", err
	}
	str := buf.String()
	return str,nil
}

func getServiceVersion(service *v1.Service) (string,error) {
	log.Debugf("getVersionOf called with service %s",service.Name)
	if len(service.Spec.Ports) == 0 {
		return "",fmt.Errorf("unable to find version - Service %s has no specified ports",service.Name)
	}
	port := service.Spec.Ports[0].Port
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d%s",service.Name,service.Namespace,port,versionPath)
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		log.Errorf("Unable to find version - failed to GET %s",url)
		return "",err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Unable to find version - failed read request body from %s",url)
			return "",err
		}
		bodyString := string(bodyBytes)
		log.Info(bodyString)
		return bodyString,nil
	} else {
		return "",fmt.Errorf("service %s %s request returned %d code",service.Name,versionPath,resp.StatusCode)
	}
}


func getServicePods(service *v1.Service,clientset *kubernetes.Clientset) (*[]v1.Pod,error) {
	log.Debugf("getPodsOf called with service %s",service.Name)
	labelSelector := labels.NewSelector()
	for labelKey, labelVal := range service.Spec.Selector {
		requirement, err := labels.NewRequirement(labelKey, selection.Equals, []string{labelVal})
		if err != nil {
			log.Error("Unable to create selector requirement")
			return nil,err
		}
		labelSelector = labelSelector.Add(*requirement)
	}
	log.Debugf("labelSelector: %s",labelSelector.String())
	runningOnlyRule := "status.phase=Running"
	options := metav1.ListOptions{
		LabelSelector:       labelSelector.String(),
		FieldSelector:       runningOnlyRule,
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), options)
	return &pods.Items,err
}

func getClientSet() (*kubernetes.Clientset,error) {
	log.Debugf("getClientSet called")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error("Unable to get cluster config")
		return nil,err
	}
	return kubernetes.NewForConfig(config)
}

func getTimeSince(startTime time.Time) time.Duration {
	log.Debugf("getTimeSince called with startTime %t",startTime)
	return time.Now().Sub(startTime)
}