package keys

import (
	"encoding/json"
	"github.com/theupdateframework/go-tuf"
	"github.com/theupdateframework/go-tuf/data"
	"github.com/theupdateframework/go-tuf/pkg/keys"
	"github.com/theupdateframework/go-tuf/sign"
	"github.com/theupdateframework/go-tuf/verify"
	"time"
)

type localKeyStore struct{

}

func (l *localKeyStore) GenKey(role string) ([]string, error) {
	return l.GenKeyWithExpires(role, data.DefaultExpires("root"))
}

func (l *localKeyStore) GenKeyWithExpires(keyRole string, expires time.Time) (keyids []string, err error) {
	signer, err := keys.GenerateEd25519Key()
	if err != nil {
		return []string{}, err
	}

	if err = l.AddPrivateKeyWithExpires(keyRole, signer, expires); err != nil {
		return []string{}, err
	}
	keyids = signer.PublicData().IDs()
	return
}

func (l *localKeyStore) AddPrivateKey(role string, signer keys.Signer) error {
	return l.AddPrivateKeyWithExpires(role, signer, data.DefaultExpires(role))
}

func (l *localKeyStore) AddPrivateKeyWithExpires(keyRole string, signer keys.Signer, expires time.Time) error {
	if !verify.ValidRole(keyRole) {
		return tuf.ErrInvalidRole{keyRole}
	}

	if !validExpires(expires) {
		return tuf.ErrInvalidExpires{expires}
	}

	if err := l.local.SaveSigner(keyRole, signer); err != nil {
		return err
	}

	if err := l.AddVerificationKeyWithExpiration(keyRole, signer.PublicData(), expires); err != nil {
		return err
	}

	return nil
}

func (l *localKeyStore) AddVerificationKey(keyRole string, pk *data.PublicKey) error {
	return l.AddVerificationKeyWithExpiration(keyRole, pk, data.DefaultExpires(keyRole))
}

func (l *localKeyStore) AddVerificationKeyWithExpiration(keyRole string, pk *data.PublicKey, expires time.Time) error {
	root, err := l.root()
	if err != nil {
		return err
	}

	role, ok := root.Roles[keyRole]
	if !ok {
		role = &data.Role{KeyIDs: []string{}, Threshold: 1}
		root.Roles[keyRole] = role
	}
	changed := false
	if role.AddKeyIDs(pk.IDs()) {
		changed = true
	}

	if root.AddKey(pk) {
		changed = true
	}

	if !changed {
		return nil
	}

	root.Expires = expires.Round(time.Second)
	if _, ok := l.versionUpdated["root.json"]; !ok {
		root.Version++
		r.versionUpdated["root.json"] = struct{}{}
	}

	return l.setMeta("root.json", root)
}

//357-395
func (l *localKeyStore) RootKeys() ([]*data.PublicKey, error) {
	root, err := l.root()
	if err != nil {
		return nil, err
	}
	role, ok := root.Roles["root"]
	if !ok {
		return nil, nil
	}

	// We might have multiple key ids that correspond to the same key, so
	// make sure we only return unique keys.
	seen := make(map[string]struct{})
	rootKeys := []*data.PublicKey{}
	for _, id := range role.KeyIDs {
		key, ok := root.Keys[id]
		if !ok {
			return nil, fmt.Errorf("tuf: invalid root metadata")
		}
		found := false
		if _, ok := seen[id]; ok {
			found = true
			break
		}
		if !found {
			for _, id := range key.IDs() {
				seen[id] = struct{}{}
			}
			rootKeys = append(rootKeys, key)
		}
	}
	return rootKeys, nil
}

func (l *localKeyStore) RevokeKey(role, id string) error {
	return l.RevokeKeyWithExpires(role, id, data.DefaultExpires("root"))
}

func (l *localKeyStore) RevokeKeyWithExpires(keyRole, id string, expires time.Time) error {
	if !verify.ValidRole(keyRole) {
		return tuf.ErrInvalidRole{keyRole}
	}

	if !validExpires(expires) {
		return tuf.ErrInvalidExpires{expires}
	}

	root, err := l.root()
	if err != nil {
		return err
	}

	key, ok := root.Keys[id]
	if !ok {
		return tuf.ErrKeyNotFound{keyRole, id}
	}

	role, ok := root.Roles[keyRole]
	if !ok {
		return tuf.ErrKeyNotFound{keyRole, id}
	}

	keyIDs := make([]string, 0, len(role.KeyIDs))

	// There may be multiple keyids that correspond to this key, so
	// filter all of them out.
	for _, keyID := range role.KeyIDs {
		if key.ContainsID(keyID) {
			continue
		}
		keyIDs = append(keyIDs, keyID)
	}
	if len(keyIDs) == len(role.KeyIDs) {
		return tuf.ErrKeyNotFound{keyRole, id}
	}
	role.KeyIDs = keyIDs

	for _, keyID := range key.IDs() {
		delete(root.Keys, keyID)
	}
	root.Roles[keyRole] = role
	root.Expires = expires.Round(time.Second)
	if _, ok := l.versionUpdated["root.json"]; !ok {
		root.Version++
		r.versionUpdated["root.json"] = struct{}{}
	}

	return l.setMeta("root.json", root)
}

//481-602
func (l *localKeyStore) Sign(roleFilename string) error {
	role := strings.TrimSuffix(roleFilename, ".json")
	if !verify.ValidRole(role) {
		return tuf.ErrInvalidRole{role}
	}

	s, err := l.SignedMeta(roleFilename)
	if err != nil {
		return err
	}

	keys, err := l.getSigningKeys(role)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return tuf.ErrInsufficientKeys{roleFilename}
	}
	for _, k := range keys {
		sign.Sign(s, k)
	}

	b, err := l.jsonMarshal(s)
	if err != nil {
		return err
	}
	r.meta[roleFilename] = b
	return l.local.SetMeta(roleFilename, b)
}

// AddOrUpdateSignature allows users to add or update a signature generated with an external tool.
// The name must be a valid metadata file name, like root.json.
func (l *localKeyStore) AddOrUpdateSignature(roleFilename string, signature data.Signature) error {
	role := strings.TrimSuffix(roleFilename, ".json")
	if !verify.ValidRole(role) {
		return tuf.ErrInvalidRole{role}
	}

	// Check key ID is in valid for the role.
	db, err := l.db()
	if err != nil {
		return err
	}
	roleData := db.GetRole(role)
	if roleData == nil {
		return tuf.ErrInvalidRole{role}
	}
	if !roleData.ValidKey(signature.KeyID) {
		return verify.ErrInvalidKey
	}

	s, err := l.SignedMeta(roleFilename)
	if err != nil {
		return err
	}

	// Add or update signature.
	signatures := make([]data.Signature, 0, len(s.Signatures)+1)
	for _, sig := range s.Signatures {
		if sig.KeyID != signature.KeyID {
			signatures = append(signatures, sig)
		}
	}
	signatures = append(signatures, signature)
	s.Signatures = signatures

	// Check signature on signed meta. Ignore threshold errors as this may not be fully
	// signed.
	if err := db.VerifySignatures(s, role); err != nil {
		if _, ok := err.(verify.ErrRoleThreshold); !ok {
			return err
		}
	}

	b, err := l.jsonMarshal(s)
	if err != nil {
		return err
	}
	r.meta[roleFilename] = b

	return l.local.SetMeta(roleFilename, b)
}

// getSigningKeys returns available signing keys.
//
// Only keys contained in the keys db are returned (i.e. local keys which have
// been revoked are omitted), except for the root role in which case all local
// keys are returned (revoked root keys still need to sign new root metadata so
// clients can verify the new root.json and update their keys db accordingly).
func (l *localKeyStore) getSigningKeys(name string) ([]keys.Signer, error) {
	signingKeys, err := l.local.GetSigners(name)
	if err != nil {
		return nil, err
	}
	if name == "root" {
		return signingKeys, nil
	}
	db, err := l.db()
	if err != nil {
		return nil, err
	}
	role := db.GetRole(name)
	if role == nil {
		return nil, nil
	}
	if len(role.KeyIDs) == 0 {
		return nil, nil
	}
	keys := make([]keys.Signer, 0, len(role.KeyIDs))
	for _, key := range signingKeys {
		for _, id := range key.PublicData().IDs() {
			if _, ok := role.KeyIDs[id]; ok {
				keys = append(keys, key)
			}
		}
	}
	return keys, nil
}

// Used to retrieve the signable portion of the metadata when using an external signing tool.
func (l *localKeyStore) SignedMeta(roleFilename string) (*data.Signed, error) {
	b, ok := l.meta[roleFilename]
	if !ok {
		return nil, tuf.ErrMissingMetadata{roleFilename}
	}
	s := &data.Signed{}
	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}
	return s, nil
}

//duplicated ones from Repo
func validExpires(expires time.Time) bool {
	return expires.Sub(time.Now()) > 0
}

func (l localKeyStore) root() (*data.Root, error) {
	rootJSON, ok := r.meta["root.json"]
	if !ok {
		return data.NewRoot(), nil
	}
	s := &data.Signed{}
	if err := json.Unmarshal(rootJSON, s); err != nil {
		return nil, err
	}
	root := &data.Root{}
	if err := json.Unmarshal(s.Signed, root); err != nil {
		return nil, err
	}
	return root, nil
}