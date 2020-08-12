package main

import (
	"flag"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
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

	stopCh :=make(chan struct{})
	defer close(stopCh)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	sharedInformers := informers.NewSharedInformerFactory(clientset,time.Minute)
	informer := sharedInformers.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj :=obj.(v1.Object)
			log.Printf("New pod added to Store:%s",mObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(v1.Object)
			nobj := newObj.(v1.Object)
			log.Printf("%s pod updated to Store:%s",oObj.GetName(),nobj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj :=obj.(v1.Object)
			log.Printf("New pod deleted to Store:%s",mObj.GetName())
		},
	})

    informer.Run(stopCh)

}
