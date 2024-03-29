package main

import (
	"path"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

var errNoData = errors.New("Found no data")

func getCredentialSetFromVault(host string) (map[string]string, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Vault client")
	}

	secret, err := client.Logical().Read(path.Join(cfg.VaultPathPrefix, host))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch secret from Vault")
	}

	if secret == nil || secret.Data == nil {
		return nil, errNoData
	}

	var (
		out  = map[string]string{}
		data map[string]interface{}
		ok   bool
	)
	if data, ok = secret.Data["data"].(map[string]interface{}); !ok {
		data = secret.Data
	}
	for k, v := range data {
		vs, ok := v.(string)
		if !ok {
			return nil, errors.Errorf("Key %q contained invalid data type %T", k, v)
		}

		out[k] = vs
	}

	return out, nil
}
