package models

import (
	"fmt"
	"sort"
	"strings"
)

type DatabaseConfig struct {
	Host             string            `json:"host"`
	Port             int               `json:"port"`
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	Database         string            `json:"database"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"` // Optional additional connection parameters mapped into the connection string
}

func quoteConfigParameter(s string) string {
	if s == "" {
		return "''"
	}
	if !strings.Contains(s, " ") {
		return s
	}
	return "'" + s + "'"
}

// Generates a connection string to be passed to sql.Open or equivalents, assuming Postgres syntax
func (c DatabaseConfig) ConnectionString() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		quoteConfigParameter(c.Host),
		c.Port,
		quoteConfigParameter(c.Username),
		quoteConfigParameter(c.Password),
		quoteConfigParameter(c.Database)))

	if _, ok := c.AdditionalParams["sslmode"]; !ok {
		b.WriteString(" sslmode=disable")
	}

	if len(c.AdditionalParams) > 0 {
		params := make([]string, 0, len(c.AdditionalParams))
		for param := range c.AdditionalParams {
			params = append(params, param)
		}

		sort.Strings(params)

		for _, param := range params {
			fmt.Fprintf(&b, " %s=%s", param, quoteConfigParameter(c.AdditionalParams[param]))
		}
	}

	return b.String()
}
