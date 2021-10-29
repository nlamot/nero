package credentials

import (
	v1 "github.com/nlamot/sofibot/pkg/k8s/core/v1"
	restclient "k8s.io/client-go/rest"
)

func List(config *restclient.Config, namespace string) ([]string, error) {
	secret, err := v1.GetSecret(config, namespace, "mmp-credentials") // mmp-credentials should be configurable
	if err != nil {
		return nil, err
	}
	return secret.GetKeys(), nil
}
