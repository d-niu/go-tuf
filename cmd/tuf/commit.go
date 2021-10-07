package main

import (
	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf/repo"
)

func init() {
	register("commit", cmdCommit, `
usage: tuf commit

Commit staged files to the repository.
`)
}

func cmdCommit(args *docopt.Args, repo *repo.Repo) error {
	return repo.Commit()
}
