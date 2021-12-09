package sealedsecrets

import (
	"fmt"

	"github.com/nlamot/nero/pkg/git"
	"github.com/nlamot/nero/pkg/gitops"
	v1 "github.com/nlamot/nero/pkg/k8s/core/v1"
)

type Manager struct {
	environment		*gitops.Environment
}

func InitSecretManager(env *gitops.Environment) (*Manager) {
	return &Manager{
		environment: env,
	}
}

func (m *Manager) WriteSecret(name string, secret *v1.Secret, changelog string) error {
	s := Sealable{secret}
	sealedSecret, err := s.Seal()
	if err != nil {
		return err
	}

	// Setup local gitops configuration
	repo, err := git.InitRepo(m.environment.Name, m.environment.SecretRepository.Location)
	if err != nil {
		return err
	}
	
	
	return nil
}
