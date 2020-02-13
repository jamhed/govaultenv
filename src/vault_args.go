package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
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

func env(name, def string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return def
}

func envb(name string, def bool) bool {
	if value, ok := os.LookupEnv(name); ok {
		if value == "true" {
			return true
		}
		return false
	}
	return def
}

const tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"

func (a *VaultArgs) Parse() *VaultArgs {
	flag.StringVar(&a.verboseLevel, "verbose", env("VERBOSE", "info"), "Set verbosity level")
	flag.StringVar(&a.kubeAuth, "kubeauth", env("KUBEAUTH", ""), "Authenticate with kubernetes, format: role@authengine")
	flag.StringVar(&a.kubeTokenPath, "kubetokenpath", env("KUBETOKENPATH", tokenPath), "Kubernetes service account token path")
	flag.StringVar(&a.vaultAddr, "addr", env("VAULT_ADDR", ""), "Vault address")
	flag.StringVar(&a.vaultToken, "token", env("VAULT_TOKEN", ""), "Vault token")
	flag.StringVar(&a.vaultTokenPath, "tokenpath", env("VAULT_TOKEN_PATH", ""), "Vault token path")
	flag.StringVar(&a.vaultPrefix, "prefix", env("VAULT_PREFIX", "VAULT_"), "Environment variable prefix")
	flag.BoolVar(&a.appendEnv, "append", envb("APPEND", true), "Append vault values to os environment")
	flag.BoolVar(&a.upperCase, "uppercase", envb("UPPERCASE", false), "Convert environment variables to upper-case")
	flag.BoolVar(&a.unwrap, "unwrap", envb("UNWRAP", false), "Unwrap token (if provided)")
	flag.BoolVar(&a.stripName, "stripname", envb("STRIPNAME", false), "Strip holding environment variable name")
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
