package v1

import (
	"bytes"
	"context"
	"encoding/base64"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

func GetSecret(config *restclient.Config, namespace string, name string) (*Secret, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return MapSecret(secret), nil
}

func (secret *Secret) Seal() ([]byte, error) {
	data, err := yaml.Marshal(secret)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("kubeseal",
		"--controller-name", "sealed-secrets", // TODO make configurable
		"--controller-namespace", "cluster-foundations", // TODO make configurable
		"--scope", "cluster-wide",
		"--format", "yaml")
	cmd.Stdin = strings.NewReader(string(data))
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	return out.Bytes(), err
}

func MapSecret(original *v1.Secret) *Secret {
	secretData := make(map[string]string)
	for k, value := range original.Data {
		secretData[k] = string(value)
	}
	secret := &Secret{
		ObjectMeta: ObjectMeta{
			Name: original.ObjectMeta.Name,
		},
		Kind:       "Secret",
		APIVersion: "v1",
		Type:       string(original.Type),
		Data:       secretData,
	}
	return secret
}

type Secret struct {
	ObjectMeta ObjectMeta        `yaml:"metadata"`
	Type       string            `yaml:"type"`
	Kind       string            `yaml:"kind"`
	APIVersion string            `yaml:"apiVersion"`
	Data       map[string]string `yaml:"data"`
}

func (secret *Secret) Update(key string, value string) {
	base64EncodedData := make([]byte, base64.StdEncoding.EncodedLen(len(value)))
	base64.StdEncoding.Encode(base64EncodedData, []byte(value))
	secret.Data[key] = string(base64EncodedData)
}

func (secret *Secret) GetKeys() ([]string) {
	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	return keys
}
