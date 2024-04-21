package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/uook3t/base/logger"
)

var (
	defaultCli = http.DefaultClient
)

const (
	HeaderAuthorization = "Authorization"
)

type RawHttpRequest struct {
	Method string
	Url    string
	Body   interface{}
	Header map[string]string

	Cli *http.Client
}

type HttpResponse struct {
	StatusCode    int         // http response status code
	Header        http.Header // http response header
	ContentLength int64       // http response content length
	Body          []byte
}

type ErrorHttpBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewRawHttpRequest(method, url string, body interface{}) *RawHttpRequest {
	return &RawHttpRequest{
		Method: method,
		Url:    url,
		Body:   body,
		Header: map[string]string{},
	}
}

func (r *RawHttpRequest) WithCli(cli *http.Client) {
	r.Cli = cli
}

func (r *RawHttpRequest) WithAccessToken(accessToken string) {
	r.Header[HeaderAuthorization] = accessToken
}

func (r *RawHttpRequest) DoSimpleHttp(ctx context.Context) (*HttpResponse, error) {
	logger.Ctx(ctx).Debugf("[DoSimpleHttp] http req start. %s %s", r.Method, r.Url)
	b, err := sonic.Marshal(r.Body)
	req, err := http.NewRequest(r.Method, r.Url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range r.Header {
		req.Header.Set(k, v)
	}

	cli := r.Cli
	if cli == nil {
		cli = defaultCli
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	httpResp := &HttpResponse{
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
		Body:          respBody,
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorHttpBody
		err = sonic.Unmarshal(respBody, &errResp)
		if err != nil {
			return nil, err
		}
		return httpResp, fmt.Errorf("do simple http failed. status: %d, body: %v", resp.StatusCode, errResp)
	}
	return httpResp, nil
}
