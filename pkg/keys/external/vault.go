package keys


type vaultClient struct {}

type defaultKey struct {}

func newVaultClient(transitPath string) (*vaultClient, error) {
	return nil, nil
}

func GenerateKey(params map[string]interface{}) (){
	path := "transitSecretEnginePathFromEnv"
	newVaultClient(path)

	//make API call using defaultKey attributes if params are empty
}

