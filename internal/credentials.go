package internal

import (
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

type Credentials struct {
	User    string
	Service string

	gitlabToken string
}

func (c *Credentials) Token() string {
	return c.gitlabToken
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
			Value:   "gitlab-tools",
			Usage:   "a system user",
			EnvVars: []string{"SERVICE"},
			Hidden:  true,
		},
	}
}

func (c *Credentials) Parse(cCtx *cli.Context) error {
	c.User = cCtx.String("system-user")
	c.Service = cCtx.String("system-service")

	secret, err := keyring.Get(c.Service, c.User)
	if err != nil {
		return err
	}

	c.gitlabToken = secret

	return nil
}
