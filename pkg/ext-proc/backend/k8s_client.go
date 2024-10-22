package backend

import (
	"context"
	"fmt"

	v1alpha1 "inference.networking.x-k8s.io/llm-instance-gateway/api/v1alpha1"
	clientset "inference.networking.x-k8s.io/llm-instance-gateway/client-go/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sClient struct {
	serverPoolName      string
	namespace           string
	kubeClientSet       *kubernetes.Clientset
	instanceGWClientSet clientset.Interface

	labelSelector string // Remove this when LLMServerPool is wired up.
}

func NewK8sClient(name, namespace string, k8sClient *kubernetes.Clientset, iGWClient clientset.Interface) K8sClient {
	return K8sClient{
		serverPoolName:      name,
		namespace:           namespace,
		kubeClientSet:       k8sClient,
		instanceGWClientSet: iGWClient,
	}
}

func (k *K8sClient) GetLLMServerPool() (*v1alpha1.LLMServerPool, error) {
	llmServerPool, err := k.instanceGWClientSet.ApiV1alpha1().LLMServerPools(k.namespace).Get(context.TODO(), k.serverPoolName, v1.GetOptions{})
	fmt.Print(llmServerPool.Name)
	if err != nil {
		fmt.Print("oh no.")
		return nil, err
	}
	return llmServerPool, nil
}

func (k *K8sClient) GetPods() {
	podList, err := k.kubeClientSet.CoreV1().Pods(k.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: k.labelSelector})
	if err != nil {
		// Handle err
	}

	for _, p := range podList.Items {
		fmt.Print(p.Name)
		// get IP and name
	}
}
