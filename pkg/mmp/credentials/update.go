package credentials

import (
	"fmt"

	restclient "k8s.io/client-go/rest"

	v1 "github.com/nlamot/sofibot/pkg/k8s/core/v1"
)

func Update(config *restclient.Config, namespace string, key string, value string) error {
	secret, err := v1.GetSecret(config, namespace, "mmp-credentials") // mmp-credentials should be configurable
	if err != nil {
		return err
	}
	secret.Update(key, value)

	sealedSecret, err := v1.SealSecret(secret)
	fmt.Println(string(sealedSecret))

	// ioutil.WriteFile("/tmp/mmp-credentials.yaml", data, 0733)

	return err
}

