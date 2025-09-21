package internal

import (
	"errors"
	"fmt"
	"net/url"
)

type GitlabPipeline struct {
	Id     int    `json:"id"`
	Ref    string `json:"ref"`
	WebUrl string `json:"web_url"`
}

func GetLatestPipeline(host, project, ref string, token Credentials) (GitlabPipeline, error) {
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
	ID       int64       `json:"id"`
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

func GetLatestJobs(host, project string, pipelineID int, token Credentials) ([]GitlabJob, error) {
	var resp []GitlabJob

	err := getResp(fmt.Sprintf("%s/api/v4/projects/%s/pipelines/%d/jobs", host, url.QueryEscape(project), pipelineID), token, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetLatestTriggeredJobs(host, project string, pipelineID int, token Credentials) ([]GitlabJob, error) {
	var resp []GitlabJob

	err := getResp(fmt.Sprintf("%s/api/v4/projects/%s/pipelines/%d/bridges", host, url.QueryEscape(project), pipelineID), token, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type RunJobResponse struct {
	Name   string `json:"name"`
	WebUrl string `json:"web_url"`
}

func RunJob(host, project string, jobID int64, token Credentials) (RunJobResponse, error) {
	var resp RunJobResponse

	err := postResp(fmt.Sprintf("%s/api/v4/projects/%s/jobs/%d/play", host, url.QueryEscape(project), jobID), token, nil, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type RetryJobResponse struct {
	Name   string `json:"name"`
	WebUrl string `json:"web_url"`
}

func RetryJob(host, project string, jobID int64, token Credentials) (RetryJobResponse, error) {
	var resp RetryJobResponse

	err := postResp(fmt.Sprintf("%s/api/v4/projects/%s/jobs/%d/retry", host, url.QueryEscape(project), jobID), token, nil, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
