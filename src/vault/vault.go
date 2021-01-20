package vault

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
)

// LoginVault : this method allows you to fetch AWS Keys securely from a Hashicorp vault instance
// Params:
// vaultAddr : vault address e.g. http://127.0.0.1:8200
// staticToken : vault login token with access to correct secret engine
// path : path to the secret KV store with relevant AWS keys present
func LoginVault(vaultAddr, staticToken, path string) (map[string]string, error) {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	client, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient})
	if err != nil {
		panic(err)
	}

	client.SetToken(staticToken)
	data, err := client.Logical().Read(path)
	if err != nil {
		panic(err)
	}

	var result map[string]string

	b, err := json.Marshal(data.Data["data"])
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &result)
	return result, err
}
