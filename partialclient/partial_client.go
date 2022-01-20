package partialclient

import (
	"encoding/json"
	"github.com/theupdateframework/go-tuf/client"
	"github.com/theupdateframework/go-tuf/data"
	"io"
	"sync"
)

type PartialClient struct {
	sync.Mutex

	rootClient  *client.Client
	localStore  partialClientLocalStore
	remoteStore partialClientRemoteStore

	valid bool

	rootVersion    uint64
	targetsVersion uint64
	targetMetas    data.TargetFiles
	//targetFiles []*pbgo.File //pbgo is dataadog specific so it's probably not waht we want
}

type partialClientLocalStore struct {
}

func (l *partialClientLocalStore) SetMeta(name string, meta json.RawMessage) error {
	return nil
}

func (l *partialClientLocalStore) GetMeta() (map[string]json.RawMessage, error) {
	return nil, nil
}

func (l *partialClientLocalStore) DeleteMeta(name string) error {
	return nil
}

func (l *partialClientLocalStore) Close() error {
	return nil
}

type partialClientRemoteStore struct {
}

func (r *partialClientRemoteStore) GetMeta(name string) (stream io.ReadCloser, size int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *partialClientRemoteStore) GetTarget(path string) (stream io.ReadCloser, size int64, err error) {
	//TODO implement me
	panic("implement me")
}

//NewPartialClient will create a partial client
func NewPartialClient(local partialClientLocalStore, remote partialClientRemoteStore, rootVersion uint64) (*PartialClient, error) {
	c := &PartialClient{
		rootClient:  client.NewClient(&local, &remote),
		localStore:  local,
		remoteStore: remote,      //do we need a custom one?
		rootVersion: rootVersion, //get the meta repo's RootsDirector's last version : check with info
	}
	return c, nil
}

//PartialUpdateAll Updates the partial client by updating roots (the partial way), root versions, and validating and updating targets
func (c *PartialClient) PartialUpdateAll() error {
	return nil
}

//FullUpdateAll Updates the partial client by updating the roots (the full way)
func (c *PartialClient) FullUpdateAll() error {
	return nil
}

//PartialUpdate Updates the partial client by updating roots (the partial way), root versions, and validating and updating targets
func (c *PartialClient) PartialUpdate(targets data.TargetFiles) error {
	return nil
}

//FullUpdate Updates the partial client by updating the roots (the full way)
func (c *PartialClient) FullUpdate(targets data.TargetFiles) error {
	return nil
}

//Targets returns the current targets (5.4.4.1 Partial Verification)
func (c *PartialClient) Targets() {

}
