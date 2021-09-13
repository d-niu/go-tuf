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
	response, err := hv.client.Logical().Read(fmt.Sprintf("%s%s", hv.transitPath, keyName))
	if err != nil {
		return nil, errors.New("Failed to read transit secret engine keys: http get failure")
	}

	keysData, hasKeys := response.Data["keys"].(map[string]interface{})
	latestVersion, hasVersion := response.Data["latest_version"].(json.Number)

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

//TODO: Parse sig data path
func (hv *hashiVaultClient) Sign(params map[string]interface{}, keyName string) (string, error) {
	signingPath := fmt.Sprintf("%s%s", hv.transitPath, keyName)
	response, err := hv.client.Logical().Write(signingPath, params)
	if err != nil {
		return "", errors.New("Failed to sign in transit secret engine: http call failed.")
	}
	sigData, ok := response.Data["signature"].(string)
	if !ok {
		return "", errors.New("Failed to sign in transit secret engine: signature corrupted.")
	}
	return sigData, nil
}

func (hv *hashiVaultClient) Verify(params map[string]interface{}, keyName string) (bool, error) {
	verifyPath := fmt.Sprintf("%s%s", hv.transitPath, keyName)
	fmt.Printf("Passing?")
	response, err := hv.client.Logical().Write(verifyPath, params)
	if err != nil {
		return false, errors.New("Failed to verify in transit secret engine: http call failed.")
	}
	valid, ok := response.Data["valid"].(bool)
	if !ok {
		return false, errors.New("Failed to verify in transit secret engine: verification response corrupted.")
	}
	return valid, nil
}

