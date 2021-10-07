package kms

import (
	"encoding/base64"
	"fmt"
	"testing"
)

const(
	transitSignPath = "/transit/sign/"
	transitVerifyPath = "/transit/verify/"
	transitKeyPath = "/transit/keys/"
	test
	testkeyname1 = "testkey1"
	dataToSign = "hello"
)

//Assumes you have a local instance of vault running.
func TestHashiVaultClient_CreateKey(t *testing.T) {
	hv, err := HVClient(transitKeyPath)
	if err != nil{
		t.Errorf("Failed to create connection to vault client: %s", err)
	}

	keytype := "ed25519"
	params := map[string]interface{}{
		"type": keytype,
	}
	err = hv.CreateKey(params, testkeyname1)
	if err != nil{
		t.Errorf("Failed to create key in vault: %s", err)
	}
}

func TestHashiVaultClient_GetPublicKey(t *testing.T) {
	hv, err := HVClient(transitKeyPath)
	if err != nil{
		t.Errorf("Failed to create connection to vault client: %s", err)
	}

	_, err = hv.GetPublicKey(testkeyname1)

	if err != nil {
		t.Errorf("Failed to read key from transit secrets engine: %s", err)
	}

	//TODO: Assert public key is data.Key type and the integrity of the publickey value
}

func TestHashiVaultClient_SignAndVerify(t *testing.T) {
	hv, err := HVClient(transitSignPath)
	if err != nil{
		t.Errorf("Failed to create connection to vault client: %s", err)
	}

	plaintext := base64.StdEncoding.EncodeToString([]byte(dataToSign))
	params := map[string]interface{}{
		"plaintext": plaintext,
	}

	signature, err := hv.Sign(params, testkeyname1)
	if err != nil {
		t.Errorf("Failed to sign in transit secret engine: %s", err)
	}

	hv2, err := HVClient(transitVerifyPath)
	if err != nil{
		t.Errorf("Failed to create connection to vault client: %s", err)
	}

	fmt.Printf("Plaintext: %s", plaintext)

	params = map[string]interface{}{
		"input" : plaintext,
		"signature": signature,
	}

	valid, err := hv2.Verify(params, testkeyname1)
	if err != nil {
		t.Errorf("Failed to verify in transit secret engine: %s", err)
	}

	if !valid{
		t.Errorf("Signature for plaintext %s failed.", dataToSign)
	}
}