package kms

import (
	"encoding/base64"
	"fmt"
	"testing"
)

const(
	transitSignPath = "/transit/sign/"
	transitKeyPath = "/transit/keys/"
	testkeyname1 = "testkey1"
	testkeyname2 = "example_edsca"
)


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

func TestHashiVaultClient_Sign(t *testing.T) {
	hv, err := HVClient(transitSignPath)
	if err != nil{
		t.Errorf("Failed to create connection to vault client: %s", err)
	}

	plaintext := base64.StdEncoding.EncodeToString([]byte("hello"))
	params := map[string]interface{}{
		"plaintext": plaintext,
	}

	ciphertext, err := hv.Sign(params, testkeyname1)
	if err != nil {
		t.Errorf("Failed to encrypt in transit secret engine %s", err)
	}

	fmt.Printf("Ciphertext is: %s", ciphertext)
}