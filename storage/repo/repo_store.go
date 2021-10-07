package repo

import (
	"encoding/json"
	"github.com/theupdateframework/go-tuf/data"
	"github.com/theupdateframework/go-tuf/pkg/keys"
	"io"
)

var TopLevelMetadata = []string{
	"root.json",
	"targets.json",
	"snapshot.json",
	"timestamp.json",
}

func signers(privateKeys []*data.PrivateKey) []keys.Signer {
	res := make([]keys.Signer, 0, len(privateKeys))
	for _, k := range privateKeys {
		signer, err := keys.GetSigner(k)
		if err != nil {
			continue
		}
		res = append(res, signer)
	}
	return res
}

type PersistedKeys struct {
	Encrypted bool            `json:"encrypted"`
	Data      json.RawMessage `json:"data"`
}

type RepoStore interface {
	// GetMeta returns a map from metadata file names (e.g. root.json) to their raw JSON payload or an error.
	GetMeta() (map[string]json.RawMessage, error)

	// SetMeta is used to update a metadata file name with a JSON payload.
	SetMeta(string, json.RawMessage) error

	// WalkStagedTargets calls targetsFn for each staged target file in paths.
	//
	// If paths is empty, all staged target files will be walked.
	WalkStagedTargets(paths []string, targetsFn TargetsWalkFunc) error

	// Commit is used to publish staged files to the repository
	Commit(bool, map[string]int, map[string]data.Hashes) error

	// GetSigners return a list of signers for a role.
	GetSigners(string) ([]keys.Signer, error)

	// SavePrivateKey adds a signer to a role.
	SaveSigner(string, keys.Signer) error

	// Clean is used to remove all staged metadata files.
	Clean() error
}

// TargetsWalkFunc is a function of a target path name and a target payload used to
// execute some function on each staged target file. For example, it may normalize path
// names and generate target file metadata with additional custom metadata.
type TargetsWalkFunc func(path string, target io.Reader) error