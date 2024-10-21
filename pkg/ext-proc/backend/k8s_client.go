package backend

import (
	"k8s.io/client-go/kubernetes"
)

type K8sClient struct {
	serverPoolName string
	clientSet      *kubernetes.Clientset
	labelSelector  string // Remove this when LLMServerPool is wired up.
}

func (k *K8sClient) GetLLMServerPool() {

}

func (k *K8sClient) GetPods(namespace string) {
	podList, err := k.clentSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: k.labelSelector})
	if err != nil {
		//Cry about it.
	}

	for p in podList.Items {
		// get IP and name
	}
}

func (k *K8sClient) GetLLMServices() {

}
