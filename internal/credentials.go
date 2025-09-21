package internal

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

type serviceKind string

func (s serviceKind) String() string {
	return string(s)
}

func (s serviceKind) Validate() (serviceKind, error) {
	switch s {
	case Gitlab, Bitbucket:
		return s, nil
	}
	return "", fmt.Errorf("unknown service: %s", s)
}

const (
	Gitlab    serviceKind = "gitlab-tools"
	Bitbucket serviceKind = "bitbucket-tools"
)

type Credentials struct {
	User    string
	Service serviceKind

	token string
}

func (c *Credentials) Token() string {
	return c.token
}

func (c *Credentials) CliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "system-user",
			Value:   "anon",
			Usage:   "a system user",
			EnvVars: []string{"USER"},
			Hidden:  true,
		},
		&cli.StringFlag{
			Name:    "system-service",
			Value:   string(Gitlab),
			Usage:   "a system user",
			EnvVars: []string{"SERVICE"},
			Hidden:  true,
		},
	}
}

func (c *Credentials) Parse(cCtx *cli.Context) (err error) {
	c.User = cCtx.String("system-user")
	c.Service, err = serviceKind(cCtx.String("system-service")).Validate()
	if err != nil {
		return err
	}

	secret, err := keyring.Get(c.Service.String(), c.User)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, keyring.ErrNotFound) {
			return keyring.Set(c.Service.String(), c.User, "")
		}
		return err
	}

	c.token = secret

	return nil
}
