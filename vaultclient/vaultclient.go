package vaultclient

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Server string // protocol, hostname and port (https://vault.foo.com:8200)
}

type VaultClient struct {
	client *api.Client
	config *VaultConfig
	token  string
}

// NewClient returns a VaultClient object or error
func NewClient(config *VaultConfig) (*VaultClient, error) {
	vc := VaultClient{}
	c, err := api.NewClient(&api.Config{Address: config.Server})
	vc.client = c
	vc.config = config
	return &vc, err
}

// TokenAuth sets the client token but doesn't check validity
func (c *VaultClient) TokenAuth(token string) {
	c.token = token
}

// AppIDAuth attempts to perform app-id authorization.
func (c *VaultClient) AppIDAuth(appid string, useridpath string) error {
	f, err := os.Open(useridpath)
	if err != nil {
		return fmt.Errorf("Error opening Vault User ID file: %v", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Error getting stat info on Vault User ID file: %v", err)
	}
	userid := make([]byte, fi.Size())
	_, err = f.Read(userid)
	if err != nil {
		return fmt.Errorf("Error reading Vault User ID file: %v", err)
	}
	req := c.client.NewRequest("POST", "/v1/auth/app-id/login")
	bodystruct := struct {
		AppID  string `json:"app_id"`
		UserID string `json:"user_id"`
	}{
		AppID:  appid,
		UserID: string(userid),
	}
	req.SetJSONBody(bodystruct)
	resp, err := c.client.RawRequest(req)
	if err != nil {
		return fmt.Errorf("Error performing auth call to Vault: %v", err)
	}
	var output interface{}
	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(&output)
	if err != nil {
		return fmt.Errorf("Error unmarshaling Vault auth response: %v", err)
	}
	body := output.(map[string]interface{})
	auth := body["auth"].(map[string]interface{})
	c.token = auth["client_token"].(string)
	return nil
}

// GetValue retrieves value at path
func (c *VaultClient) GetValue(path string) (interface{}, error) {
	c.client.SetToken(c.token)
	lc := c.client.Logical()
	s, err := lc.Read(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading secret from Vault: %v: %v", path, err)
	}
	return s.Data["value"], nil
}

// WriteValue writes value=data at path
func (c *VaultClient) WriteValue(path string, data []byte) error {
	c.client.SetToken(c.token)
	lc := c.client.Logical()
	_, err := lc.Write(path, map[string]interface{}{"value": data})
	return err
}
