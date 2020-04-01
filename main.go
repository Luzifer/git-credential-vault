package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/v2/env"
	"github.com/Luzifer/rconfig/v2"
)

const (
	actionGet   = "get"
	actionStore = "store"
	actionErase = "erase"
)

var (
	cfg = struct {
		LogLevel        string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VaultAddr       string `flag:"vault-addr" default:"" description:"Vault Address to connect to" validate:"nonzero"`
		VaultPathPrefix string `flag:"vault-path-prefix" default:"" description:"Prefix to prepend to hostname" validate:"nonzero"`
		VaultToken      string `flag:"vault-token" default:"" description:"Token to use in Vault connection" validate:"nonzero"`
		VersionAndExit  bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

// Instead of init() function not to fail tests by arg parser
func initApp() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("git-credential-vault %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	initApp()

	var action string
	if len(rconfig.Args()) == 2 {
		action = rconfig.Args()[1]
	}

	switch action {

	case actionGet:
		handleGetAction()

	case actionErase, actionStore:
		// Defined action but unsupported by this provider
		return

	default:
		log.Fatalf("Unsupported action %q", action)

	}
}

func handleGetAction() {
	values, err := readInput(os.Stdin)
	if err != nil {
		log.WithError(err).Fatal("Unable to parse input values")
	}

	if proto := values["protocol"]; proto != "https" {
		log.Debugf("Unsupported protocol %q", proto)
		return
	}

	vaultValues, err := getCredentialSetFromVault(values["host"])
	if err != nil {
		log.WithError(err).Fatal("Unable to retrieve values from Vault")
	}

	for k, v := range vaultValues {
		values[k] = v
	}

	fmt.Fprintln(os.Stdout, strings.Join(env.MapToList(values), "\n"))
}

func readInput(input io.Reader) (map[string]string, error) {
	var (
		lines   []string
		scanner = bufio.NewScanner(input)
	)

	for scanner.Scan() {
		var text = strings.TrimSpace(scanner.Text())

		if text == "" {
			break
		}

		lines = append(lines, text)
	}

	return env.ListToMap(lines), errors.Wrap(scanner.Err(), "Unable to read input")
}
