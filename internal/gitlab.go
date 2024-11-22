package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/urfave/cli/v2"
)

type GitlabInfo struct {
	GitlabHost  string
	RefName     string
	ProjectName string
}

func (g *GitlabInfo) CliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "gitlab-host",
			Value: "",
			Usage: "gitlab host",
		},
		&cli.StringFlag{
			Name:  "ref",
			Value: "master",
			Usage: "git ref",
		},
		&cli.StringFlag{
			Name:  "gitlab-project",
			Value: "",
			Usage: "gitlab project name",
		},
		&cli.BoolFlag{
			Name:  "use-git-cli",
			Value: true,
			Usage: "using git cli",
		},
	}
}

func (g *GitlabInfo) Parse(cCtx *cli.Context) error {
	g.GitlabHost = cCtx.String("gitlab-host")
	g.RefName = cCtx.String("ref")
	g.ProjectName = cCtx.String("gitlab-project")

	useGitCli := cCtx.Bool("use-git-cli")

	if useGitCli {
		ref, err := GetCurrentGitRef()
		if err != nil {
			return err
		}

		host, project, err := GetCurrentGitProjectName()
		if err != nil {
			return err
		}

		g.GitlabHost = host
		g.RefName = ref
		g.ProjectName = project
	}

	return nil
}

var httpClient = http.DefaultClient

func getResp[T any](url, token string, dst *T) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s, response: %s", resp.Status, string(b))
	}

	err = json.Unmarshal(b, dst)
	if err != nil {
		return err
	}

	return nil
}

type GitlabPipeline struct {
	Id     int    `json:"id"`
	Ref    string `json:"ref"`
	WebUrl string `json:"web_url"`
}

func GetLatestPipeline(host, project, ref, token string) (GitlabPipeline, error) {
	var resp []GitlabPipeline

	err := getResp(fmt.Sprintf("%s/api/v4/projects/%s/pipelines?ref=%s&per_page=%d", host, url.QueryEscape(project), ref, 1), token, &resp)
	if err != nil {
		return GitlabPipeline{}, err
	}

	if len(resp) == 0 {
		return GitlabPipeline{}, errors.New("not found a gitlab pipeline")
	}

	return resp[0], nil
}

type GitlabJob struct {
	Status   string      `json:"status"`
	Stage    string      `json:"stage"`
	Name     string      `json:"name"`
	Coverage interface{} `json:"coverage"`
	Duration *float64    `json:"duration"`
	User     struct {
		Username string `json:"username"`
		Name     string `json:"name"`
	} `json:"user"`
	WebUrl string `json:"web_url"`
}

func GetLatestJobs(host, project string, pipelineID int, token string) ([]GitlabJob, error) {
	var resp []GitlabJob

	err := getResp(fmt.Sprintf("%s/api/v4/projects/%s/pipelines/%d/jobs", host, url.QueryEscape(project), pipelineID), token, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetLatestTriggeredJobs(host, project string, pipelineID int, token string) ([]GitlabJob, error) {
	var resp []GitlabJob

	err := getResp(fmt.Sprintf("%s/api/v4/projects/%s/pipelines/%d/bridges", host, url.QueryEscape(project), pipelineID), token, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
