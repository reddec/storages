package rest

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)

func NewClient(baseURL string) *restClient {
	return NewClientContext(baseURL, context.Background(), DefaultTimeout)
}

func NewClientContext(baseURL string, ctx context.Context, timeout time.Duration) *restClient {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &restClient{
		ctx:     ctx,
		timeout: timeout,
		baseURL: baseURL,
	}
}

type restClient struct {
	ctx     context.Context
	timeout time.Duration
	baseURL string
}

func (r *restClient) Put(key []byte, data []byte) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+base64.StdEncoding.EncodeToString(key), bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "rest: post key, prepare request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "rest: post key, execute request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("rest: post key, unexpected status %v %v", res.StatusCode, res.Status)
	}
	return nil
}

func (r *restClient) Close() error { return nil }

func (r *restClient) Get(key []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+base64.StdEncoding.EncodeToString(key), nil)
	if err != nil {
		return nil, errors.Wrap(err, "rest: get key, prepare request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "rest: get key, execute request")
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, os.ErrNotExist
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("rest: get key, unexpected status %v %v", res.StatusCode, res.Status)
	}
	return ioutil.ReadAll(res.Body)
}

func (r *restClient) Del(key []byte) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, r.baseURL+base64.StdEncoding.EncodeToString(key), nil)
	if err != nil {
		return errors.Wrap(err, "rest: delete key, prepare request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "rest: delete key, execute request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("rest: delete key, unexpected status %v %v", res.StatusCode, res.Status)
	}
	return nil
}

func (r *restClient) Keys(handler func(key []byte) error) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL, nil)
	if err != nil {
		return errors.Wrap(err, "rest: list keys, prepare request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "rest: list keys, execute request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("rest: list keys, %v", res.Status)
	}
	reader := bufio.NewScanner(res.Body)
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if len(line) == 0 {
			continue
		}
		key, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return errors.Wrapf(err, "rest: decode key %v", line)
		}
		err = handler(key)
		if err != nil {
			return err
		}
	}
	return nil
}
