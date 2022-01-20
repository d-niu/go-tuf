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
	store       partialClientLocalStore
	repo        *tuf.Repo
	local       partialClientLocalStore
	remote      partialClientRemoteStore
	expiredTime time.Time
	keyIDs      map[string][]string
}

//Test errors that should be detected upon initializing any partial client
func (s *PartialClientSuite) TestInit(c *C) {

}

//Test partial (as defined in Uptane) updating all targets
func (s *PartialClientSuite) TestPartialUpdateAll() {

}

//Test full (as defined in TUF) updating all targets
func (s *PartialClientSuite) TestFullUpdateAll() {

}

//Test partial (as defined in Uptane) updating some targets
func (s *PartialClientSuite) TestPartialUpdate() {

}

//Test full (as defined in TUF) updating some targets
func (s *PartialClientSuite) TestFullUpdate() {

}

//TODO: test valid/invalid rotations, missing/incorrect target files
