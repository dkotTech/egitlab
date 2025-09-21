package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

var httpClient = http.DefaultClient

type HttpSettings struct {
	WithFakeHttpTransport bool
}

func (c *HttpSettings) CliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "fake-http-transport",
			Value: false,
			Usage: "enable fake http transport with requests logs",
		},
	}
}

func (c *HttpSettings) Parse(cCtx *cli.Context) (err error) {
	c.WithFakeHttpTransport = cCtx.Bool("fake-http-transport")

	if c.WithFakeHttpTransport {
		httpClient.Transport = &FakeTransport{}
	}

	return nil
}

func getResp[T any](url string, token Credentials, dst *T) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	switch token.Service {
	case Gitlab:
		req.Header.Set("PRIVATE-TOKEN", token.Token())
	case Bitbucket:
		req.Header.Set("Authorization", "Basic "+token.Token())
	}

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

func postResp[T any](url string, token Credentials, body io.Reader, dst *T) error {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	switch token.Service {
	case Gitlab:
		req.Header.Set("PRIVATE-TOKEN", token.Token())
	case Bitbucket:
		req.Header.Set("Authorization", "Basic "+token.Token())
	}

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

type FakeTransport struct{}

func (t *FakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf(">>> %s %s\n", req.Method, req.URL.String())

	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		fmt.Printf("Body: %s\n", string(body))
	}

	return &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("{}")),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}
