package partialclient

import (
	"github.com/theupdateframework/go-tuf/client"
	"github.com/theupdateframework/go-tuf/data"
	"sync"
)

type PartialClient struct {
	sync.Mutex

	rootClient  *client.Client
	localStore  client.LocalStore
	remoteStore client.RemoteStore //*partialClientRemoteStore

	valid bool

	rootVersion    uint64
	targetsVersion uint64
	targetMetas    data.TargetFiles
	//targetFiles []*pbgo.File //pbgo is dataadog specific so it's probably not waht we want
}

//NewPartialClient will create a partial client
func NewPartialClient(local client.LocalStore, remote client.RemoteStore, rootVersion uint64) (*PartialClient, error) {
	c := &PartialClient{
		rootClient:  client.NewClient(local, remote),
		localStore:  local,
		remoteStore: remote,      //do we need a custom one?
		rootVersion: rootVersion, //get the meta repo's RootsDirector's last version : check with info
	}
	return c, nil
}

//Update the partial client by updating roots, root versions, and validating and updating targets
func (c *PartialClient) Update() error {
	return nil
}

//Targets returns the current targets (5.4.4.1 Partial Verification)
func (c *PartialClient) Targets() {
	return
}
