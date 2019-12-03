package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

// Default is `-s -w -X main.version={{.Version}} -X main.commit={{.ShortCommit}} -X main.date={{.Date}} -X main.builtBy=goreleaser`.
var version string
var commit string
var date string
var builtBy string

type VaultArgs struct {
	verboseLevel   string
	kubeAuth       string
	kubeTokenPath  string
	vaultAddr      string
	vaultToken     string
	vaultTokenPath string
	vaultPrefix    string
	appendEnv      bool
	upperCase      bool
	unwrap         bool
	stripName      bool
	args           []string
}

func NewArgs() *VaultArgs {
	return new(VaultArgs)
}

func (a *VaultArgs) Parse() *VaultArgs {
	flag.StringVar(&a.verboseLevel, "verbose", "info", "Set verbosity level")
	flag.StringVar(&a.kubeAuth, "kubeauth", "", "Authenticate with kubernetes, format: role@authengine")
	flag.StringVar(&a.kubeTokenPath, "kubetokenpath", "/var/run/secrets/kubernetes.io/serviceaccount/token", "Kubernetes service account token path")
	flag.StringVar(&a.vaultAddr, "addr", os.Getenv("VAULT_ADDR"), "Vault address")
	flag.StringVar(&a.vaultToken, "token", "", "Vault token")
	flag.StringVar(&a.vaultTokenPath, "tokenpath", "", "Vault token path")
	flag.StringVar(&a.vaultPrefix, "prefix", "VAULT_", "Environment variable prefix")
	flag.BoolVar(&a.appendEnv, "append", true, "Append vault vaulues to os environment")
	flag.BoolVar(&a.upperCase, "uppercase", true, "Convert environment variables to upper-case")
	flag.BoolVar(&a.unwrap, "unwrap", false, "Unwrap token (if provided)")
	flag.BoolVar(&a.stripName, "stripname", false, "Strip holding environment variable name")
	flag.Parse()
	a.args = flag.Args()
	return a
}

func (a *VaultArgs) LogLevel() *VaultArgs {
	if a.verboseLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if a.verboseLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if a.verboseLevel == "error" {
		log.SetLevel(log.ErrorLevel)
	} else if a.verboseLevel == "fatal" {
		log.SetLevel(log.FatalLevel)
	} else if a.verboseLevel == "panic" {
		log.SetLevel(log.PanicLevel)
	}
	return a
}

func (a *VaultArgs) Validate() *VaultArgs {
	if len(a.args) == 0 {
		fmt.Printf("Usage: govaultenv [-h] command [arguments]\n")
		fmt.Printf("version:%s commit:%s build by:%s date:%s\n", version, commit, builtBy, date)
		os.Exit(1)
	}
	return a
}
