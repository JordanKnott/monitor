package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type MonitorClient struct {
	httpClient graphql.Doer
	endpoint   string
	token      string
}

func NewClient(endpoint string, httpClient graphql.Doer, token string) graphql.Client {
	if httpClient == nil || httpClient == (*http.Client)(nil) {
		httpClient = http.DefaultClient
	}
	return &MonitorClient{httpClient, endpoint, token}
}

func (c *MonitorClient) createPostRequest(req *graphql.Request) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		c.endpoint,
		bytes.NewReader(body))
	httpReq.Header.Add("Authorization", c.token)
	if err != nil {
		return nil, err
	}

	return httpReq, nil
}

func (c MonitorClient) MakeRequest(
	ctx context.Context,
	req *graphql.Request,
	resp *graphql.Response,
) error {

	var httpReq *http.Request
	var err error
	httpReq, err = c.createPostRequest(req)

	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if ctx != nil {
		httpReq = httpReq.WithContext(ctx)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		var respBody []byte
		respBody, err = io.ReadAll(httpResp.Body)
		if err != nil {
			respBody = []byte(fmt.Sprintf("<unreadable: %v>", err))
		}
		return fmt.Errorf("returned error %v: %s", httpResp.Status, respBody)
	}

	err = json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		return err
	}
	if len(resp.Errors) > 0 {
		return resp.Errors
	}
	return nil
}
