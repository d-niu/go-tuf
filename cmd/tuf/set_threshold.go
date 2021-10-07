package main

import (
	"fmt"
	"github.com/theupdateframework/go-tuf/repo"
	"strconv"

	"github.com/flynn/go-docopt"
)

func init() {
	register("set-threshold", cmdSetThreshold, `
usage: tuf set-threshold <role> <threshold>

Set the threshold for a role.  
`)
}

func cmdSetThreshold(args *docopt.Args, repo *repo.Repo) error {
	role := args.String["<role>"]
	thresholdStr := args.String["<threshold>"]
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		return err
	}

	if err := repo.SetThreshold(role, threshold); err != nil {
		return err
	}

	fmt.Println("Set ", role, "threshold to", threshold)
	return nil
}
