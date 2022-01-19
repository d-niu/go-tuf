package partialclient

import (
	"github.com/theupdateframework/go-tuf"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

//Mirroring the setup for the full verification client
type PartialClientSuite struct {
	store       tuf.LocalStore
	repo        *tuf.Repo
	local       tuf.LocalStore
	remote      *fakeRemoteStore //make one
	expiredTime time.Time
	keyIDs      map[string][]string
}

type fakeRemoteStore struct {
}

func newFakeRemoteStore() *fakeRemoteStore {
	return &fakeRemoteStore{}
}

func (s *PartialClientSuite) TestInit(c *C) {

}
