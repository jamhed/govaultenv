package main

import (
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

type VaultMockApi struct {
	client *api.Client
}

func (v *VaultMockApi) Get(path string) VaultData {
	return VaultData{"key1": "value1", "key2": "value2"}
}

func (v *VaultMockApi) SetClient(client *api.Client) {
	v.client = client
}

func (v *VaultMockApi) Client() *api.Client {
	return v.client
}

func NewVaultTest(upperCase, stripName bool) *Vault {
	v := NewVault("vault_addr", upperCase, stripName)
	v.api = new(VaultMockApi)
	return v
}

func Test_Secret(t *testing.T) {
	v := NewVaultTest(true, false)
	v.ParseValue("TEST", "path")
	assert.Equal(t, len(v.env), 2)
	assert.Contains(t, v.env, "TEST_KEY1=value1")
	assert.Contains(t, v.env, "TEST_KEY2=value2")
}

func Test_Secret_Key(t *testing.T) {
	v := NewVaultTest(true, false)
	v.ParseValue("TEST", "path#key1")
	assert.Equal(t, len(v.env), 1)
	assert.Equal(t, v.env[0], "TEST=value1")
}

func Test_Secret_With_Strip_Name(t *testing.T) {
	v := NewVaultTest(true, true)
	v.ParseValue("TEST", "path")
	assert.Equal(t, len(v.env), 2)
	assert.Contains(t, v.env, "KEY1=value1")
	assert.Contains(t, v.env, "KEY2=value2")
}

func Test_Secret_With_No_Upper_Case(t *testing.T) {
	v := NewVaultTest(false, false)
	v.ParseValue("TEST", "path")
	assert.Equal(t, len(v.env), 2)
	assert.Contains(t, v.env, "TEST_key1=value1")
	assert.Contains(t, v.env, "TEST_key2=value2")
}

func Test_Secret_File(t *testing.T) {
	v := NewVaultTest(true, false)
	v.ParseValue("TEST", "path#key1:filepath")
	filepath := "filepath"
	assert.Equal(t, len(v.env), 0)
	assert.Equal(t, len(v.file), 1)
	assert.Equal(t, v.file[0], filepath)
	assert.FileExists(t, filepath)
	content, err := ioutil.ReadFile(filepath)
	assert.Nil(t, err, "file content read")
	assert.Equal(t, string(content), "value1")
	v.Cleanup()
	_, err = os.Stat(filepath)
	assert.True(t, os.IsNotExist(err))
}
