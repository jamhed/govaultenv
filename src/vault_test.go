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
