package sealedsecrets

import (
	"bytes"
	"os/exec"
	"strings"

	v1 "github.com/nlamot/nero/pkg/k8s/core/v1"
	"gopkg.in/yaml.v2"
)

type Sealable struct {
	*v1.Secret
}


func (secret *Sealable) Seal() ([]byte, error) {
	data, err := yaml.Marshal(secret)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("kubeseal",
		"--controller-name", "sealed-secrets", // TODO make configurable
		"--controller-namespace", "cluster-foundations", // TODO make configurable
		"--scope", "cluster-wide",
		"--format", "yaml")
	cmd.Stdin = strings.NewReader(string(data))
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	return out.Bytes(), err
}
