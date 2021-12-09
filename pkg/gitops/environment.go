package gitops

import (
	v1 "github.com/nlamot/nero/pkg/k8s/core/v1"
	restclient "k8s.io/client-go/rest"
)

type Environment struct {
	Name            	string
	StateRepository		GitOpsStateRepository
	SecretRepository   	SecretRepository
	K8sClientConfig 	*restclient.Config // Should be a repo
	SecretManager		SecretManager
}

type GitRepository struct {
	Location  string
	Branch    string
	Directory string
}

func (e *Environment) UpdateCredentials(key string, value string) error {
	secret, err := v1.GetSecret(e.K8sClientConfig, e.Name, "mmp-credentials") // mmp-credentials should be configurable
	if err != nil {
		return err
	}
	secret.Update(key, value)

	err = e.SecretManager.WriteSecret("", secret, "test")
	if err != nil {
		return err
	}

	return err
}

func (e *Environment) ListCredentials() ([]string, error) {
	secret, err := v1.GetSecret(e.K8sClientConfig, e.Name, "mmp-credentials") // mmp-credentials should be configurable
	if err != nil {
		return nil, err
	}
	return secret.GetKeys(), nil
}
