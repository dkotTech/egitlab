package internal

import "github.com/urfave/cli/v2"

type ServiceInfo struct {
	Host          string
	RefName       string
	ProjectName   string
	LastCommitMsg string
	HeadBranch    string
}

func (g *ServiceInfo) CliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Usage: "service host",
		},
		&cli.StringFlag{
			Name:  "ref",
			Value: "master",
			Usage: "git ref",
		},
		&cli.StringFlag{
			Name:  "project",
			Usage: "project name",
		},
		&cli.StringFlag{
			Name:  "last-commit-msg",
			Usage: "last commit message",
		},
		&cli.StringFlag{
			Name:  "head-branch",
			Usage: "head branch",
		},
		&cli.BoolFlag{
			Name:  "use-git-cli",
			Value: true,
			Usage: "using git cli",
		},
		&cli.BoolFlag{
			Name:  "use-git-cli-head-request",
			Value: false,
			Usage: "make a git head request, its may request a ssh key password",
		},
	}
}

func (g *ServiceInfo) Parse(cCtx *cli.Context) error {
	g.Host = cCtx.String("host")
	g.RefName = cCtx.String("ref")
	g.ProjectName = cCtx.String("project")
	g.LastCommitMsg = cCtx.String("last-commit-msg")
	g.HeadBranch = cCtx.String("head-branch")

	useGitCli := cCtx.Bool("use-git-cli")
	useGitHeadRequest := cCtx.Bool("use-git-cli-head-request")

	if useGitCli {
		var err error

		g.RefName, err = GetCurrentGitRef()
		if err != nil {
			return err
		}

		g.Host, g.ProjectName, err = GetCurrentGitProjectName()
		if err != nil {
			return err
		}

		g.LastCommitMsg, err = GetLastCommitMessage()
		if err != nil {
			return err
		}

		if useGitHeadRequest {
			g.HeadBranch, err = GetHeadBranch()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
