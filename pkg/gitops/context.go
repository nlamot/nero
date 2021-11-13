package gitops

type ContextRepository interface {
	Get(name string) Context
	List() []Context			// TODO
	Add(context Context)		// TODO
}

type Context struct {
	Name string
	Environments []Environment
}

