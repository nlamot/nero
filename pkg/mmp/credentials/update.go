package credentials

import (
	"fmt"

	restclient "k8s.io/client-go/rest"

	"github.com/nlamot/nero/pkg/git"
	v1 "github.com/nlamot/nero/pkg/k8s/core/v1"
)

func Update(clientConfig *restclient.Config, namespace string, key string, value string) error {
	secret, err := v1.GetSecret(clientConfig, namespace, "mmp-credentials") // mmp-credentials should be configurable
	if err != nil {
		return err
	}
	secret.Update(key, value)

	sealedSecret, err := secret.Seal()
	if err != nil {
		return err
	}

	// Setup local gitops configuration
	repo, err := git.InitRepo("clt-config", "git@bitbucket.org:sofico/clt-config.git", "master")
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
			Path: fmt.Sprintf("env/%s/general-config/mmp-credentials.yaml", namespace),
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
	err = repo.Checkout("master", true)
	if err != nil {
		return err
	}

	return err
}
