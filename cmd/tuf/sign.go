package main

import (
	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf/repo"
)

func init() {
	register("sign", cmdSign, `
usage: tuf sign <metadata>

Sign a role's metadata file.

Signs the given role's staged metadata file with all keys present in the 'keys'
directory for that role.
`)
}

func cmdSign(args *docopt.Args, repo *repo.Repo) error {
	return repo.Sign(args.String["<metadata>"])
}
