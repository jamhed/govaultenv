package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
)

func getToken(vaultToken, vaultTokenPath string) string {
	if len(vaultToken) > 0 {
		return vaultToken
	}
	token := os.Getenv("VAULT_TOKEN")
	if len(token) > 0 {
		log.Debugf("Using token from environment variable VAULT_TOKEN")
		return token
	}
	var tokenPath string
	if len(vaultTokenPath) > 0 {
		tokenPath = vaultTokenPath
	} else {
		user, err := user.Current()
		if err == nil {
			tokenPath = path.Join(user.HomeDir, ".vault-token")
		} else {
			log.Errorf("Can't find user's home, error:%s", err)
			os.Exit(1)
		}
	}
	log.Debugf("Reading vault token from:%s", tokenPath)
	byteToken, err := ioutil.ReadFile(tokenPath)
	if err == nil {
		return string(byteToken)
	} else {
		log.Errorf("Can't read token from path:%s, error:%s", tokenPath, err)
		os.Exit(1)
	}
	return ""
}

func main() {
	args := NewArgs().Parse().Validate().LogLevel()

	log.Debugf("govaultenv, version:%s commit:%s build by:%s date:%s\n", version, commit, builtBy, date)

	v := NewVault(args.vaultAddr, args.upperCase, args.stripName).Connect()

	if len(args.kubeAuth) > 0 {
		re := regexp.MustCompile(`^(.+?)@(.+)$`).FindStringSubmatch(args.kubeAuth)
		if len(re) == 3 {
			v.kubeTokenPath = args.kubeTokenPath
			token := v.KubeAuth(re[1], re[2])
			v.SetToken(false, token)
		} else {
			log.Errorf("Can't parse kubeauth option value:%s", args.kubeAuth)
			os.Exit(1)
		}
	} else {
		v.SetToken(args.unwrap, getToken(args.vaultToken, args.vaultTokenPath))
	}

	reEnv := regexp.MustCompile(fmt.Sprintf("^%s(.+?)=(.+)$", args.vaultPrefix))
	for _, env := range os.Environ() {
		if val := reEnv.FindStringSubmatch(env); len(val) == 3 {
			envName, envValue := val[1], val[2]
			if args.vaultPrefix == "VAULT_" && (envName == "ADDR" || envName == "TOKEN") {
				log.Debugf("Skipping env variable name:%s%s", args.vaultPrefix, envName)
				continue
			}
			log.Debugf("Parsing secret:%s from env:%s%s", envValue, args.vaultPrefix, envName)
			v.ParseSecretPath(envName, envValue)
		}
	}

	cmd := exec.Command(args.args[0], args.args[1:]...)
	if args.appendEnv {
		cmd.Env = append(os.Environ(), v.env...)
	} else {
		cmd.Env = v.env
	}
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	log.Debugf("Starting command:%s, arguments:%s", args.args[0], args.args[1:])
	err := cmd.Run()

	v.Cleanup()

	if err != nil {
		log.Fatal(err)
	}
}
