package kms

import(
	"encoding/json"
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/theupdateframework/go-tuf/data"
	"os"
)

func HVClient(keyResourceID string) (KMSClient, error) {
	hv, err := newHashiVaultClient(keyResourceID)
	if err != nil {
		return nil, err
	}
	return hv, nil
}

type hashiVaultClient struct {
	client			*vault.Client
	transitPath 	string
}

//TODO: Parameterize getting vault address and tokens.
//TODO: Add error messages to the package
func newHashiVaultClient(transitPath string) (*hashiVaultClient, error) {
	address := os.Getenv("VAULT_ADDR")
	if address == "" {
		return nil, errors.New("No vault address set")
	}

	token :=  os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, errors.New("No dev access token")
	}

	//Instantiate new Vault Client
	client, err := vault.NewClient(&vault.Config{
		Address: address,
	})
	if err != nil{
		return nil, err
	}
	client.SetToken(token)

	//do we need transit path if we're calling hv client on a method-needed basis?
	hvClient := &hashiVaultClient{
		client: client,
		transitPath: transitPath,
	}

	return hvClient, nil
}

//CreateKey: returns the result of a request to generate a ed25519 key in Vault (data.Key Type and Scheme not parameterized atm).
func (hv *hashiVaultClient) CreateKey(params map[string]interface{}, keyName string) (error) {
	//TODO: Create with Role (key role or creator role).
	newkeypath := fmt.Sprintf("%s%s", hv.transitPath, keyName)
	_, err := hv.client.Logical().Write(newkeypath, params)
	if err != nil {
		return err
	}
	return  nil
}

func (hv *hashiVaultClient) GetPublicKey(keyName string) (*data.Key, error){
	secret, err := hv.client.Logical().Read(fmt.Sprintf("%s%s", hv.transitPath, keyName))
	if err != nil {
		return nil, errors.New("Failed to read transit secret engine keys: http get failure")
	}

	keysData, hasKeys := secret.Data["keys"].(map[string]interface{})
	latestVersion, hasVersion := secret.Data["latest_version"].(json.Number)

	if !hasKeys || !hasVersion {
		return nil, errors.New("Failed to read transit secret engine keys: corrupted initial response.")
	}

	keyData, ok := keysData[string(latestVersion)].(map[string]interface{})
	if !ok {
		return nil, errors.New("Failed to read transit secret engine keys: key specific data corrupted.")
	}

	publicKey, ok := keyData["public_key"].(string)
	if !ok {
		return nil, errors.New("Failed to read transit secret engine keys: public key data corrupted.")
	}

	//TODO: dynamically retrieve these from the key read function or createkey parameters
	publickey := data.Key{
		Type: data.KeyTypeEd25519,
		Scheme: data.KeySchemeEd25519,
		Algorithms: data.KeyAlgorithms,
		Value: json.RawMessage(publicKey),
	}

	return &publickey, nil
}

func (hv *hashiVaultClient) Sign(params map[string]interface{}, keyName string) (string, error) {
	signingPath := fmt.Sprintf("%s%s", hv.transitPath, keyName)
	secret, err := hv.client.Logical().Write(signingPath, params)
	if err != nil {
		return "", err
	}
	cipherData, ok := secret.Data["signature"].(string)
	if !ok {
		return "", errors.New("Failed to encrypt in transit secret engine: signature corrupted.")
	}
	return cipherData, nil
}

func (hv *hashiVaultClient) Verify()  {
	return
}

