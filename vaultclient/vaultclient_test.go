package vaultclient

import (
	"log"
	"testing"
)

var tconfig = VaultConfig{
	Server:     "https://vault-prod.shave.io:8200",
	AppID:      "acquisition-development",
	UserIDPath: "testdata/userid.txt",
}

func TestVaultAppIDAuth(t *testing.T) {
	vc, err := NewClient(&tconfig)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	err = vc.AppIDAuth()
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
}

func TestVaultGetValue(t *testing.T) {
	vc, err := NewClient(&tconfig)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	err = vc.AppIDAuth()
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
	d, err := vc.GetValue("secret/development/acquisition/test_value")
	if err != nil {
		log.Fatalf("Error getting value: %v", err)
	}
	log.Printf("Got value: %v", d.(string))
}
