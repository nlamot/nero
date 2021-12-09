package git

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitRepo struct {
	url        string
	localPath  string
	mainBranch string
	repo       *git.Repository
	auth       *ssh.PublicKeysCallback
	worktree   *git.Worktree
}

type GitRepoChange struct {
	Path string
	Data []byte
}

func InitRepo(name string, url string, mainBranch string) (*GitRepo, error) {
	fmt.Printf("Setting up repo %s", name)
	g := &GitRepo{
		localPath:  fmt.Sprintf("/tmp/%s", name), // clt-config
		url:        url,                          // "git@bitbucket.org:sofico/clt-config.git"
		mainBranch: mainBranch,                   // master
	}
	var err error
	g.auth, err = ssh.NewSSHAgentAuth("git")
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(g.localPath); os.IsNotExist(err) {
		g.repo, err = git.PlainClone(g.localPath, false, &git.CloneOptions{
			Auth: g.auth,
			URL:  g.url,
		})
		if err != nil {
			return nil, err
		}
	} else {
		g.repo, err = git.PlainOpenWithOptions(g.localPath, &git.PlainOpenOptions{})
		if err != nil {
			return nil, err
		}
	}
	g.worktree, err = g.repo.Worktree()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GitRepo) CreateFeatureBranch(name string) error {
	fmt.Printf("Creating feature branch %s", name)
	branch := fmt.Sprintf("refs/heads/%s", name)

	if g.HasBranch(branch) {
		err := g.Checkout(name, true)
		if err != nil {
			return err
		}
	} else {
		g.Checkout("master", true)
		err := g.worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branch),
			Create: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GitRepo) StageChanges(changes []GitRepoChange) error {
	fmt.Println("Staging changes")
	for _, change := range changes {
		if change.Data != nil {
			path := fmt.Sprintf("%s/%s", g.localPath, change.Path)
			// Update content
			err := ioutil.WriteFile(path, change.Data, 0633)
			if err != nil {
				return err
			}
			g.worktree.Add(path)
		}
	}
	return nil
}

func (g *GitRepo) CommitAndPushStagedChanges(message string) (*string, error) {
	fmt.Printf("Committing change %s", message)
	g.worktree.Commit(fmt.Sprintf("[NERO] %s", message), &git.CommitOptions{})
	var output bytes.Buffer
	err := g.repo.Push(&git.PushOptions{
		Auth:     g.auth,
		Progress: &output,
	})
	fmt.Println(output.String())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (g *GitRepo) Checkout(branch string, force bool) error {
	fmt.Printf("Checkout branch %s", branch)
	err := g.worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  force,
	})
	if err != nil {
		return err
	}
	err = g.worktree.Pull(&git.PullOptions{
		Force: force,
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func (g *GitRepo) CheckoutMain(force bool) error {
	return g.Checkout(g.mainBranch, force)
}

func (g *GitRepo) HasBranch(name string) bool {
	branch := plumbing.ReferenceName(name)
	refs, _ := g.repo.Branches()
	exists := false
	refs.ForEach(func(r *plumbing.Reference) error {
		if r.Name() == branch {
			exists = true
			return nil
		}
		return nil
	})
	return exists
}
