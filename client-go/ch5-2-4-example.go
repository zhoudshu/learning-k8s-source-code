package main

import (
	"flag"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 配置 k8s 集群外 kubeconfig 配置文件
	var kubeconfig *string
		kubeconfig = flag.String("kubeconfig", "./admin.conf", "absolute path to the kubeconfig file")
	flag.Parse()

	//在 kubeconfig 中使用当前上下文环境，config 获取支持 url 和 path 方式
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 根据指定的 config 创建一个新的 dynamicClient
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取指定 namespace 中的 Pod 列表信息
	gvr :=schema.GroupVersionResource{Version: "v1", Resource:"pods"}
	unstrunctObj, err := dynamicClient.Resource(gvr).Namespace(apiv1.NamespaceDefault).List(metav1.ListOptions{Limit:500})
	if err != nil {
		panic(err)
	}
	podList := &corev1.PodList{}
	err =runtime.DefaultUnstructuredConverter.FromUnstructured(unstrunctObj.UnstructuredContent(), podList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nThere are %d pods in namespaces %s\n", len(podList.Items), apiv1.NamespaceDefault)
	for _, pod := range podList.Items {
		fmt.Printf("Name: %s, Status: %s, CreateTime: %s\n", pod.ObjectMeta.Name, pod.Status.Phase, pod.ObjectMeta.CreationTimestamp)
	}

}
