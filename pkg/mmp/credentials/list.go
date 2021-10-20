package credentials

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

func List(config *restclient.Config) ([]string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	secret, err := clientset.CoreV1().Secrets("default").Get(context.Background(), "mmp-credentials", v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	
	return keys, nil
}