package keys

import "github.com/theupdateframework/go-tuf/pkg/keys"

type ExternalKey interface {
	GenerateKey(params map[string]interface{}) (keyReference string, err error)

	//Repo handles updating signatures wrt roles and files
	SignMessage(message []byte) ([]byte, error)

	GetVerifier(keyReference string) (keys.Verifier, error)

	Revoke(keyReference string) error
}



