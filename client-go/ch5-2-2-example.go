package main

import (
	"flag"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
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

	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err.Error())
	}

	result := &corev1.PodList{}
	err = restClient.Get().
		Namespace("default").
		Resource("pods").
		VersionedParams(&metav1.ListOptions{Limit:500}, scheme.ParameterCodec).
		Do().
		Into(result)

	if err != nil {
		panic(err.Error())
	}

	for _, d := range result.Items {
		fmt.Printf("Namespace:%v \t Name:%v \t Status:%+v\n", d.Namespace,d.Name,d.Status.Phase)
	}
}
