package main

import (
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
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

	// 根据指定的 config 创建一个新的 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取指定 namespace 中的 Pod 列表信息
	namespace := "default"
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{Limit:500})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nThere are %d pods in namespaces %s\n", len(pods.Items), namespace)
	for _, pod := range pods.Items {
		fmt.Printf("Name: %s, Status: %s, CreateTime: %s\n", pod.ObjectMeta.Name, pod.Status.Phase, pod.ObjectMeta.CreationTimestamp)
	}

}

func prettyPrint(maps map[string]interface{}) {
	lens := 0
	for k, _ := range maps {
		if lens <= len(k) {
			lens = len(k)
		}
	}
	for key, values := range maps {
		spaces := lens - len(key)
		v := ""
		for i := 0; i < spaces; i++ {
			v += " "
		}
		fmt.Printf("%s: %s%v\n", key, v, values)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
