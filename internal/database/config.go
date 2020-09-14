package database

import (
	"net/url"
)

type Config struct {
	Name     string `envconfig:"DATASERVICE_DB_NAME"`
	Host     string `envconfig:"DATASERVICE_DB_HOST"`
	Port     string `envconfig:"DATASERVICE_DB_PORT"`
	User     string `envconfig:"DATASERVICE_DB_USER"`
	Password string `envconfig:"DATASERVICE_DB_PASSWORD"`
	SSLMode  string `envconfig:"DATASERVICE_DB_SSLMODE"`
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	host := c.Host
	if v := c.Port; v != "" {
		host = host + ":" + v
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   host,
		Path:   c.Name,
	}

	q := u.Query()
	if v := c.SSLMode; v != "" {
		q.Add("sslmode", v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
