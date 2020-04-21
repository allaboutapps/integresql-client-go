package integresql

import "github.com/allaboutapps/integresql-client-go/pkg/util"

type ClientConfig struct {
	BaseURL    string
	APIVersion string
}

func DefaultClientConfigFromEnv() ClientConfig {
	return ClientConfig{
		BaseURL:    util.GetEnv("INTEGRESQL_CLIENT_BASE_URL", "http://integresql:5000/api"),
		APIVersion: util.GetEnv("INTEGRESQL_CLIENT_API_VERSION", "v1"),
	}
}
