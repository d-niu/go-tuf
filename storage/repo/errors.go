package repo

import (
	"errors"
	"fmt"
)

var (
	ErrNewRepository  = errors.New("tuf: repository not yet committed")
)

type ErrFileNotFound struct {
	Path string
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("tuf: file not found %s", e.Path)
}

type ErrPassphraseRequired struct {
	Role string
}

func (e ErrPassphraseRequired) Error() string {
	return fmt.Sprintf("tuf: a passphrase is required to access the encrypted %s keys file", e.Role)
}
