package kms

import (
	"github.com/theupdateframework/go-tuf/data"
)

//Client is an interface because a client performs the actions listed below.
type KMSClient interface {
	//TODO: add specific parameters you'd need + no name for keys
	CreateKey(params map[string]interface{}, keyname string) (err error)

	//fetch signer and verifier for the key <- delegations

	//TODO: fix ciphertext naming
	Sign(params map[string]interface{}, keyname string) (signature string, err error)

	Verify(params map[string]interface{}, keyname string) (verified bool, err error)

	GetPublicKey(path string) (*data.Key, error)

	//DeleteKey(params map[string]interface{}) (bool, error)
}


type CryptoManager interface {
	//Right now we have:
		//Signer map -> no maps
		//Verifier map -> no maps

	//Parse root metadata file and get keys in mem

	//Ideally: initialize tuf client with details for the Signer backend
		//Right now: client.db <- client specific map



}