package main

import (
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

type VaultApi struct {
	client *api.Client
}

func (v *VaultApi) Get(path string) VaultData {
	re, err := v.client.Logical().Read(path)
	if err != nil {
		log.Errorf("Secret:%s, %s", path, err)
		return make(VaultData)
	}
	if re == nil {
		log.Warnf("No data for secret, path:%s", path)
		return make(VaultData)
	}
	if data, ok := re.Data["data"].(map[string]interface{}); ok {
		return data
	}
	return re.Data
}

func (v *VaultApi) SetClient(client *api.Client) {
	v.client = client
}

func (v *VaultApi) Client() *api.Client {
	return v.client
}
