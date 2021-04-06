package gitops

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/fluxcd/flux2/pkg/manifestgen/install"
)

type InstallOptions struct {
}

func Install(options InstallOptions) error {
	tmpDir, err := ioutil.TempDir("", "kubectl-aks")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	opts := install.MakeDefaultOptions()
	m, err := install.Generate(opts)
	if err != nil {
		return err
	}

	path, err := m.WriteFile(tmpDir)
	if err != nil {
		return err
	}

	if err := apply(path); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func apply(path string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	command := exec.CommandContext(ctx, "/bin/sh", "-c", fmt.Sprintf("kubectl apply -f %s --cache-dir=/tmp --dry-run=client", path))
	output, err := command.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}
