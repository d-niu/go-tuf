package main

import (
	"github.com/theupdateframework/go-tuf/repo"
	"log"

	"github.com/flynn/go-docopt"
)

func init() {
	register("regenerate", cmdRegenerate, `
usage: tuf regenerate [--consistent-snapshot=false]

Recreate the targets metadata file. Important: Not supported yet

Alternatively, passphrases can be set via environment variables in the
form of TUF_{{ROLE}}_PASSPHRASE
`)
}

func cmdRegenerate(args *docopt.Args, repo *repo.Repo) error {
	// TODO: implement this
	log.Println("Not supported yet")
	return nil
}
