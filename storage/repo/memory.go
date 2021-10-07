package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/theupdateframework/go-tuf/data"
	"github.com/theupdateframework/go-tuf/pkg/keys"
)

func MemoryStore(meta map[string]json.RawMessage, files map[string][]byte) RepoStore {
	if meta == nil {
		meta = make(map[string]json.RawMessage)
	}
	return &memoryStore{
		meta:       meta,
		stagedMeta: make(map[string]json.RawMessage),
		files:      files,
		signers:    make(map[string][]keys.Signer),
	}
}

type memoryStore struct {
	meta       map[string]json.RawMessage
	stagedMeta map[string]json.RawMessage
	files      map[string][]byte
	signers    map[string][]keys.Signer
}

func (m *memoryStore) GetMeta() (map[string]json.RawMessage, error) {
	meta := make(map[string]json.RawMessage, len(m.meta)+len(m.stagedMeta))
	for key, value := range m.meta {
		meta[key] = value
	}
	for key, value := range m.stagedMeta {
		meta[key] = value
	}
	return meta, nil
}

func (m *memoryStore) SetMeta(name string, meta json.RawMessage) error {
	m.stagedMeta[name] = meta
	return nil
}

func (m *memoryStore) WalkStagedTargets(paths []string, targetsFn TargetsWalkFunc) error {
	if len(paths) == 0 {
		for path, data := range m.files {
			if err := targetsFn(path, bytes.NewReader(data)); err != nil {
				return err
			}
		}
		return nil
	}

	for _, path := range paths {
		data, ok := m.files[path]
		if !ok {
			return fmt.Errorf("tuf: file not found %s", path)
		}
		if err := targetsFn(path, bytes.NewReader(data)); err != nil {
			return err
		}
	}
	return nil
}

func (m *memoryStore) Commit(consistentSnapshot bool, versions map[string]int, hashes map[string]data.Hashes) error {
	for name, meta := range m.stagedMeta {
		paths := computeMetadataPaths(consistentSnapshot, name, versions)
		for _, path := range paths {
			m.meta[path] = meta
		}
	}
	return nil
}

func (m *memoryStore) GetSigners(role string) ([]keys.Signer, error) {
	return m.signers[role], nil
}

func (m *memoryStore) SaveSigner(role string, signer keys.Signer) error {
	m.signers[role] = append(m.signers[role], signer)
	return nil
}

func (m *memoryStore) Clean() error {
	return nil
}

//mock or fake service (simulates https server - local) -> asserts same number of tokens etc
//separate the logic for stuff

//separate PR for tests
	//google form for feedback for people who use the tool
//separate PR for metrics
