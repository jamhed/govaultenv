package main

import (
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

type VaultApi struct {
	client *api.Client
}

func path2v2(secretPath string) string {
	base, secret := path.Split(secretPath)
	return fmt.Sprintf("%sdata/%s", base, secret)
}

func (v *VaultApi) Get(path string) VaultData {
	pathv2 := path2v2(path)
	re, err := v.client.Logical().Read(pathv2)
	if err != nil {
		log.Errorf("Secret:%s, %s", pathv2, err)
		return make(VaultData)
	}
	if re != nil {
		if data, ok := re.Data["data"].(map[string]interface{}); ok {
			return data
		}
		log.Warnf("No data key for secret:%s", pathv2)
	}

	log.Warnf("Fall back to version 1, secret:%s", path)
	re, err = v.client.Logical().Read(path)
	if err != nil {
		log.Errorf("Secret:%s, %s", path, err)
		return make(VaultData)
	}
	if re != nil {
		return re.Data
	}
	log.Warnf("No either v2 or v1 data for secret:%s", path)
	return make(VaultData)
}

func (v *VaultApi) SetClient(client *api.Client) {
	v.client = client
}

func (v *VaultApi) Client() *api.Client {
	return v.client
}
