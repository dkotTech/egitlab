package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type CreatePRRequest struct {
	Title  string `json:"title"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"destination"`
	CloseSourceBranch bool `json:"close_source_branch"`
}

type CreatePRResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"self"`
	} `json:"links"`
}

func CreatePR(host, project string, title, description, sourceBranch, destinationBranch string, token Credentials) (CreatePRResponse, error) {
	buf := bytes.NewBuffer(nil)

	err := json.NewEncoder(buf).Encode(&CreatePRRequest{
		Title: title,
		Source: struct {
			Branch struct {
				Name string `json:"name"`
			} `json:"branch"`
		}{
			Branch: struct {
				Name string `json:"name"`
			}{Name: sourceBranch},
		},
		Destination: struct {
			Branch struct {
				Name string `json:"name"`
			} `json:"branch"`
		}{
			Branch: struct {
				Name string `json:"name"`
			}{Name: destinationBranch},
		},
		CloseSourceBranch: true,
	})
	if err != nil {
		return CreatePRResponse{}, err
	}

	var res CreatePRResponse
	err = postResp(fmt.Sprintf("%s/repositories/%s/pullrequests", host, project), token, buf, &res)
	if err != nil {
		return CreatePRResponse{}, err
	}

	return res, nil
}
