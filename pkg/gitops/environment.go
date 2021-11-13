package gitops

import (
	"fmt"

	"github.com/nlamot/nero/pkg/git"

	v1 "github.com/nlamot/nero/pkg/k8s/core/v1"
	restclient "k8s.io/client-go/rest"
)

type Environment struct {
	Name            string
	GitRepository   SourceReference
	K8sClientConfig *restclient.Config // Should be a repo
}

type SourceReference struct {
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

	sealedSecret, err := secret.Seal()
	if err != nil {
		return err
	}

	// Setup local gitops configuration
	repo, err := git.InitRepo(e.Name, e.GitRepository.Location)
	if err != nil {
		return err
	}
	// Return to master state
	err = repo.Checkout(e.GitRepository.Branch, true)
	if err != nil {
		return err
	}

	// Create branch to work on
	err = repo.CreateFeatureBranch("nero-test")
	if err != nil {
		return err
	}

	// Update mmp-credentials & get URL to create PR
	err = repo.StageChanges([]git.GitRepoChange{
		{
			Path: fmt.Sprintf("%s/general-config/mmp-credentials.yaml", e.GitRepository.Directory),
			Data: sealedSecret,
		},
	})
	if err != nil {
		return err
	}
	_, err = repo.CommitAndPushStagedChanges("Test")
	if err != nil {
		return err
	}
	// Return to master state
	err = repo.Checkout(e.GitRepository.Branch, true)
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
