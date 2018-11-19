package client

import (
	"client-go/kubernetes"
	restclient "client-go/rest"
	"client-go/tools/clientcmd"
	"fmt"
	log "github.com/sirupsen/logrus"
	"k8s.io/kubernetes/panda/model"
	"sync"
)

var locker sync.RWMutex
var mapKubernetesClientSet = make(map[string]*kubernetes.Clientset)

func GetKuberntesClientSet(label string) (*kubernetes.Clientset, error) {
	// 从全局变量中读取 k8s client
	clientSet, ok := ReadMap(label)
	if ok {
		return clientSet, nil
	}

	// 如果当前的字典中没有,那么从数据库中获取
	InitK8sClientSetFromDb()  // admin.config 全部和数据库同步
	clientSet, ok = ReadMap(label)
	if !ok {
		return nil, fmt.Errorf(" can't found kubernetes client in cluster (%s)", label)
	}
	return clientSet, nil
}

func ReadMap(label string) (*kubernetes.Clientset, bool) {
	locker.RLocker()
	clientSet, ok := mapKubernetesClientSet[label]
	locker.RUnlock()
	return clientSet, ok
}

func WriteMap(label string, clientSet *kubernetes.Clientset) {
	locker.Lock()
	mapKubernetesClientSet[label] = clientSet
	locker.Unlock()
}

func InitK8sClientSetFromDb() {
	log.Info("Init InitK8sClientSetByDb......")
	result := model.GetKubeAdmin("MANAGE")
	for _, v := range result {
		err := AddK8sClientToMap(v.Label, []byte(v.Context))
		if err != nil {
			continue
		}
	}
}

// AddK8sClientToMap
func AddK8sClientToMap(label string, context []byte) error {
	clientSet, err := NewK8sClientSetWithContext(context)
	if err != nil {
		log.Errorf("AddK8sClientToMap,label:%s ,err:%s", label, err.Error())
		return err
	}
	WriteMap(label, clientSet)
	log.Infof("AddK8sClientToMap success,label:%s", label)
	return nil
}

func NewK8sClientSetWithContext(context []byte) (*kubernetes.Clientset, error) {
	restConfig, err := GetRestConfigWithContext(context)
	if err != nil {
		log.Errorf("init k8s ClientConfig failed, %v", err)
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Errorf("init k8s clientSet failed, %v", err)
		return nil, err
	}

	return clientSet, nil
}

func GetRestConfigWithContext(context []byte) (*restclient.Config, error) {
	config, err := clientcmd.Load(context)
	if err != nil {
		log.Errorf("k8s load config failed, %v", err)
		return nil, err
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})
	return clientConfig.ClientConfig()
}
