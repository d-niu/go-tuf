package kms

import (
	"github.com/theupdateframework/go-tuf/data"
)

//Client is an interface because a client performs the actions listed below.
//Defer individual client structs to specific kms.
type KMSClient interface {
	//TODO: Handle different key types. Further, change publickey back to data.Publickey type to fit into the implementation
	CreateKey(params map[string]interface{}, keyname string) (err error)

	Sign()

	Verify()

	GetPublicKey(path string) (*data.Key, error)
}
