package config

import (
	"os"

	"github-graph-drawer/log"
)

var (
	_config = config{}
)

func init() {
	_config = config{
		Port: getEnv("API_PORT", "8080"),
		Mailer: mailerConf{
			host:     getEnv("MAILER_HOST", ""),
			port:     getEnv("MAILER_PORT", ""),
			user:     getEnv("MAILER_USER", ""),
			password: getEnv("MAILER_PASSWORD", ""),
		},
		DbUri: getEnv("DB_URI", ""),
	}
}

type config struct {
	Port   string
	Mailer mailerConf
	DbUri  string
}

type mailerConf struct {
	host     string
	port     string
	user     string
	password string
}

// Config returns the API's config :)
func Config() config {
	return _config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Warningf("The \"%s\" variable is not set. Defaulting to \"%s\".\n", key, defaultValue)
		value = defaultValue
	}
	return value
}
