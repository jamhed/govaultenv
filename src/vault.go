package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

type VaultApiInterface interface {
	Get(string) VaultData
	SetClient(*api.Client)
	Client() *api.Client
}

type Vault struct {
	api           VaultApiInterface
	addr          string
	kubeTokenPath string
	upperCase     bool
	stripName     bool
	env           []string
	envFilter     map[string]bool
	file          []string
}

type VaultData map[string]interface{}

var reEnvKey = regexp.MustCompile(`^(.+?)#(.+)$`)
var reFileKey = regexp.MustCompile(`^(.+?)#(.+?):(.+)$`)

func NewVault(addr string, upperCase, stripName bool) *Vault {
	v := new(Vault)
	v.upperCase = upperCase
	v.stripName = stripName
	v.addr = addr
	v.file = make([]string, 0)
	v.env = make([]string, 0)
	v.envFilter = make(map[string]bool)
	v.api = new(VaultApi)
	return v
}

func (v *Vault) Connect() *Vault {
	log.Debugf("Connecting to vault addr:%s", v.addr)
	c, err := api.NewClient(&api.Config{Address: v.addr})
	if err != nil {
		log.Errorf("Failed to connect to vault addr:%s, error:%s", v.addr, err)
		os.Exit(1)
	}
	v.api.SetClient(c)
	return v
}

func (v *Vault) SetToken(unwrap bool, token string) {
	if !unwrap {
		v.api.Client().SetToken(token)
		return
	}
	re, err := v.api.Client().Logical().Unwrap(token)
	if err != nil {
		log.Errorf("Can't unwrap token:%s, error:%s", token, err)
		os.Exit(1)
	}
	v.api.Client().SetToken(re.Auth.ClientToken)
}

func (v *Vault) Filter(envs []string) (re []string) {
	for _, env := range envs {
		if !v.envFilter[env] {
			re = append(re, env)
		}
	}
	return
}

func (v *Vault) GetValue(path string, key string) interface{} {
	re := v.api.Get(path)
	if value, ok := re[key]; ok {
		return value
	}
	if data, ok := re["data"].(map[string]interface{}); ok {
		if data[key] == nil {
			log.Warnf("Empty key:%s value in path:%s", key, path)
			return ""
		}
		return data[key]
	} else {
		log.Warnf("No data for key:%s value in path:%s", key, path)
		return ""
	}
}

func (v *Vault) KubeAuth(role, path string) string {
	jwt, err := ioutil.ReadFile(v.kubeTokenPath)
	if err != nil {
		log.Errorf("Can't read jwt token at %s", v.kubeTokenPath)
		os.Exit(1)
	}
	re, err := v.api.Client().Logical().Write("auth/"+path+"/login", map[string]interface{}{"role": role, "jwt": string(jwt)})
	if err != nil {
		log.Errorf("Can't authenticate jwt token path:%s, role:%s, error:%s", path, role, err)
		os.Exit(1)
	}
	return re.Auth.ClientToken
}

func (v *Vault) makeEnvVarName(envName, key string) (name string) {
	if v.stripName {
		name = key
	} else {
		name = fmt.Sprintf("%s_%s", envName, key)
	}
	if v.upperCase {
		name = strings.ToUpper(name)
	}
	return name
}

func (v *Vault) writeFile(name, content string) {
	if err := ioutil.WriteFile(name, []byte(content), 0644); err != nil {
		log.Errorf("file:%s, %s", name, err)
	} else {
		log.Debugf("Wrote secret file:%s", name)
		v.file = append(v.file, name)
	}
}

func (v *Vault) handleSecretFile(envName, secretPath, keyName, localFileName string) {
	secretValue := v.GetValue(secretPath, keyName)
	switch value := secretValue.(type) {
	case string:
		v.writeFile(localFileName, value)
	default:
		log.Errorf("Can't write file value of type:%s", value)
	}
}

func (v *Vault) handleSecretKey(envName, secretPath, keyName string) {
	log.Debugf("Set env variable name:%s", envName)
	v.env = append(v.env, fmt.Sprintf("%s=%s", envName, v.GetValue(secretPath, keyName)))
}

func (v *Vault) handleSecret(envName, envValue string) {
	for secretKey, secretValue := range v.api.Get(envValue) {
		name := v.makeEnvVarName(envName, secretKey)
		log.Debugf("Set env variable name:%s", name)
		v.env = append(v.env, fmt.Sprintf("%s=%s", name, secretValue))
	}
}

func (v *Vault) ParseSecretPath(envName, envValue string) {
	if valParts := reFileKey.FindStringSubmatch(envValue); len(valParts) == 4 {
		v.handleSecretFile(envName, valParts[1], valParts[2], valParts[3])
	} else if valParts := reEnvKey.FindStringSubmatch(envValue); len(valParts) == 3 {
		v.handleSecretKey(envName, valParts[1], valParts[2])
	} else {
		v.handleSecret(envName, envValue)
	}
}

func (v *Vault) Cleanup() {
	for _, file := range v.file {
		log.Debugf("Remove secret file:%s", file)
		if err := os.Remove(file); err != nil {
			log.Errorf("file:%s, %s", file, err)
		}
	}
}
