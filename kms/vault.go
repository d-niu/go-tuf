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
	client	*vault.Client
	transitSecretEnginePath string
}

//TODO: what about using const inputs?
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

	hvClient := &hashiVaultClient{
		client: client,
		transitSecretEnginePath: transitPath,
	}

	return hvClient, nil
}

//CreateKey: returns the result of a request to generate a ed25519 key in Vault (data.Key Type and Scheme not parameterized atm).
func (hv *hashiVaultClient) CreateKey(params map[string]interface{}, keyname string) (error) {
	//TODO: Create with Role (key role or creator role). Investigate if we need this functoinality
	//equivalent of cli 'vault write -f keymgmt/key/example-key type '
	newkeypath := fmt.Sprintf("%s%s", hv.transitSecretEnginePath, keyname)
	_, err := hv.client.Logical().Write(newkeypath, params)
	if err != nil {
		return err
	}
	return  nil
}

func (hv *hashiVaultClient) GetPublicKey(newkeypath string) (*data.Key, error){
	secret, err := hv.client.Logical().Read(newkeypath)
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

func (hv *hashiVaultClient) Sign() {
	return
}


func (hv *hashiVaultClient) Verify()  {
	return
}

