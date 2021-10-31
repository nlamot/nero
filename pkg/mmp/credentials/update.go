package credentials

import (
	"fmt"
	"io/ioutil"
	"os"

	restclient "k8s.io/client-go/rest"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	v1 "github.com/nlamot/sofibot/pkg/k8s/core/v1"
)

func Update(clientConfig *restclient.Config, namespace string, key string, value string) error {
	secret, err := v1.GetSecret(clientConfig, namespace, "mmp-credentials") // mmp-credentials should be configurable
	if err != nil {
		return err
	}
	secret.Update(key, value)

	sealedSecret, err := v1.SealSecret(secret)
	if err != nil {
		return err
	}

	var repo *git.Repository
	publicKeys, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		return err
	}
	if _, err := os.Stat("/tmp/clt-config"); os.IsNotExist(err) {
		repo, err = git.PlainClone("/tmp/clt-config", false, &git.CloneOptions{
			Auth:     publicKeys,
			URL:      "git@bitbucket.org:sofico/clt-config.git",
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	} else {
		repo, err = git.PlainOpenWithOptions("/tmp/clt-config", &git.PlainOpenOptions{
			DetectDotGit: true,
		})
		if err != nil {
			return err
		}
	}
	branch := fmt.Sprintf("refs/heads/%s", "nero-test")
	b := plumbing.ReferenceName(branch)
	refs, _ := repo.Branches()
	worktree, _ := repo.Worktree()
	exists := false
	refs.ForEach(func(r *plumbing.Reference) error {
		if r.Name() == b {
			exists = true
			return nil
		}
		return nil
	})
	if exists {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: b,
			Force:  true,
		})
		if err != nil {
			return err
		}
		worktree.Pull(&git.PullOptions{
			Force: true,
		})
	} else {
		worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", "master")),
			Force:  true,
		})
		fmt.Println("checkout")
		worktree.Pull(&git.PullOptions{
			Force: true,
		})
		fmt.Println("pull")
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: b,
			Create: true,
		})
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	fmt.Println("worktree")
	
	fmt.Println("checkout nero-test")
	ioutil.WriteFile(fmt.Sprintf("/tmp/clt-config/env/%s/general-config/mmp-credentials.yaml", namespace), sealedSecret, 0633)
	fmt.Println("writefile")
	worktree.Add(fmt.Sprintf("env/%s/general-config/mmp-credentials.yaml", namespace))
	fmt.Println("add")
	worktree.Commit("Nero test", &git.CommitOptions{})
	fmt.Println("commit")
	err = repo.Push(&git.PushOptions{
		Auth: publicKeys,
	})
	if err != nil {
		return err
	}
	fmt.Println("push")
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", "master")),
	})
	if err != nil {
		return err
	}
	fmt.Println("checkout master")

	return err
}
