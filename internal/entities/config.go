package entities

import (
	"encoding/base64"
	"fmt"
)

type Config struct {
	Grafana grafana `json:"grafana"`
}

type grafana struct {
	Auth auth   `json:"auth"`
	Host string `json:"host"`
}

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a auth) BasicHeader() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(a.Username+":"+a.Password)))
}

type ApiKey string

func (k ApiKey) BearerHeader() string {
	return fmt.Sprintf("Bearer %s", k)
}

type Session string

func (s Session) CookieHeader() string {
	return fmt.Sprintf("Cookie: %s", s)
}
