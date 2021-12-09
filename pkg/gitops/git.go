package gitops

import (
	"net/url"

	"github.com/nlamot/nero/pkg/git"
)

type GitOpsStateRepository struct {
	*git.GitRepo
}

func (g *GitOpsStateRepository) RequestChanges(branch string, changes []git.GitRepoChange, changelog string) (*url.URL, error){
	// Checkout main state
	err := g.CheckoutMain(true)
	if err != nil {
		return nil, err
	}

	// Create branch to work on
	err = g.CreateFeatureBranch(branch)
	if err != nil {
		return nil, err
	}

	// Update mmp-credentials & get URL to create PR
	err = g.StageChanges(changes)
	if err != nil {
		return nil, err
	}
	_, err = g.CommitAndPushStagedChanges(changelog)
	if err != nil {
		return nil, err
	}
	// Return to main state
	err = g.CheckoutMain(true)
	if err != nil {
		return nil, err
	}
	return nil, nil	// TODO extract PR url
}