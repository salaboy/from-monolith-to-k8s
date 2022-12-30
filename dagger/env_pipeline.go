package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	if len(os.Args) < 3 {
		log.Fatalf("Incorrect arguments: expected repo, KubeCfgPath")
	}

	err := run(ctx, os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context, repo, kubeCfgPath string) error {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}

	defer c.Close()

	kubeCfgFilePath := c.Host().Directory(filepath.Dir(kubeCfgPath)).File(filepath.Base(kubeCfgPath))

	repoDir := c.Git(repo).Branch("main").Tree()

	_, err = c.Container().From("quay.io/roboll/helmfile:helm3-v0.135.0").
		WithWorkdir("/app").
		WithMountedDirectory(".", repoDir).
		WithFile(".kube/config", kubeCfgFilePath).
		WithEnvVariable("FOO", time.Now().String()).
		WithExec([]string{"sh", "-c", `
set -ex
helmfile repos
helmfile template > final.yaml
helmfile sync
`}).ExitCode(ctx)

	return err
}
