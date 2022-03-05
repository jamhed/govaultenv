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
	return fmt.Sprintf("%s/data/%s", base, secret)
}

func (v *VaultApi) Get(path string) VaultData {
	pathv2 := path2v2(path)
	re, err := v.client.Logical().Read(pathv2)
	if err != nil {
		log.Errorf("Secret:%s, %s", path, err)
		return make(VaultData)
	}
	if re != nil {
		if data, ok := re.Data["data"].(map[string]interface{}); ok {
			return data
		}
		log.Warnf("%s %s", pathv2, re)
		log.Warnf("No data key for secret:%s", path)
	}
	re, err = v.client.Logical().Read(path)
	if err != nil {
		log.Errorf("Secret:%s, %s", path, err)
		return make(VaultData)
	}
	if re != nil {
		if len(re.Data) == 0 {
			log.Warnf("No data for secret:%s", path)
		}
		return re.Data
	}
	log.Warnf("No data for secret:%s", path)
	return make(VaultData)
}

func (v *VaultApi) SetClient(client *api.Client) {
	v.client = client
}

func (v *VaultApi) Client() *api.Client {
	return v.client
}
