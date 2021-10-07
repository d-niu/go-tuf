package main

import (
	"fmt"
	"github.com/flynn/go-docopt"
	"github.com/theupdateframework/go-tuf/repo"
)

func init() {
	register("get-threshold", cmdGetThreshold, `
usage: tuf get-threshold <role>

Gets the threshold for a role.  
`)
}

func cmdGetThreshold(args *docopt.Args, repo *repo.Repo) error {
	role := args.String["<role>"]

	threshold, err := repo.GetThreshold(role)
	if err != nil {
		return err
	}

	fmt.Println("Got", role, "threshold", threshold)
	return nil
}
