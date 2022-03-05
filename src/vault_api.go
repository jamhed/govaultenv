package main

import (
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

type VaultApi struct {
	client *api.Client
}

func (v *VaultApi) Get(path string) VaultData {
	mountPath, v2, err := isKVv2(path, v.client)
	if err != nil {
		log.Errorf("%s", err.Error())
		return make(VaultData)
	}
	if v2 {
		path = addPrefixToKVPath(path, mountPath, "data")
	}
	var versionParam map[string]string
	secret, err := kvReadRequest(v.client, path, versionParam)
	if err != nil {
		log.Errorf("path:%s, %s", path, err.Error())
		if secret != nil {
			return secret.Data
		}
		return make(VaultData)
	}
	if secret == nil {
		log.Warnf("No value found at %s", path)
		return make(VaultData)
	}
	if v2 {
		if data, ok := secret.Data["data"].(map[string]interface{}); ok && data != nil {
			return data
		} else {
			log.Errorf("No data found at %s", path)
			return make(VaultData)
		}
	} else {
		return secret.Data
	}
}

func (v *VaultApi) SetClient(client *api.Client) {
	v.client = client
}

func (v *VaultApi) Client() *api.Client {
	return v.client
}
