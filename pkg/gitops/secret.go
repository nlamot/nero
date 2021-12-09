package gitops

import (
	v1 "github.com/nlamot/nero/pkg/k8s/core/v1"
)

type SecretManager interface {
	WriteSecret(name string, secret *v1.Secret, changelog string) error
}

type SecretRepository interface {
}
