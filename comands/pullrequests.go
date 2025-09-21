package comands

import (
	"encoding/json"
	"fmt"

	"github.com/dkotTech/egitlab/internal"
	"github.com/urfave/cli/v2"
)

func NewPullRequestsCommand() *cli.Command {
	commandFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "job",
			Value: "create",
		},
	}

	var credentials internal.Credentials
	var serviceInfo internal.ServiceInfo
	var httpSettings internal.HttpSettings

	commandFlags = append(commandFlags, credentials.CliFlags()...)
	commandFlags = append(commandFlags, serviceInfo.CliFlags()...)
	commandFlags = append(commandFlags, httpSettings.CliFlags()...)

	return &cli.Command{
		Name:    "pullrequests",
		Aliases: []string{"pr"},
		Usage:   "work with pull requests",
		Flags:   commandFlags,
		Action: func(cCtx *cli.Context) error {
			err := credentials.Parse(cCtx)
			if err != nil {
				return err
			}

			err = serviceInfo.Parse(cCtx)
			if err != nil {
				return err
			}

			err = httpSettings.Parse(cCtx)
			if err != nil {
				return err
			}

			fields := make(map[string]string)
			fieldsData := cCtx.Args().First()
			if fieldsData != "" {
				err = json.Unmarshal([]byte(fieldsData), &fields)
				if err != nil {
					return err
				}
			}

			app := pullRequestsApp{
				job:         cCtx.String("job"),
				credentials: credentials,
				serviceInfo: serviceInfo,
				fields:      fields,
			}

			return app.doJob()()
		},
	}
}

type pullRequestsApp struct {
	job string

	credentials internal.Credentials
	serviceInfo internal.ServiceInfo

	fields map[string]string
}

func (p pullRequestsApp) doJob() func() error {
	switch p.job {
	case "create":
		return p.CreatePR
	}
	return func() error {
		return nil
	}
}

func (p pullRequestsApp) CreatePR() error {
	res, err := internal.CreatePR(p.serviceInfo.Host, p.serviceInfo.ProjectName,
		getOrDefault(p.fields["title"], p.serviceInfo.LastCommitMsg),
		p.fields["description"],
		getOrDefault(p.fields["source_branch"], p.serviceInfo.RefName),
		getOrDefault(p.fields["destination_branch"], p.serviceInfo.HeadBranch),
		p.credentials,
	)
	if err != nil {
		return err
	}

	fmt.Println("pull request created: ", res.Links.Self.Href)

	return nil
}

func getOrDefault[T comparable](v T, dv T) T {
	var def T
	if def == v {
		return dv
	}
	return v
}
