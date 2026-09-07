package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	Token      string
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (c *Client) SetToken(token string) {
	c.Token = token
}

func (c *Client) ClearToken() {
	c.Token = ""
}

func (c *Client) Get(
	address string,
	path string,
	result any,
) error {
	return c.do(
		http.MethodGet,
		address,
		path,
		nil,
		result,
		true,
	)
}

func (c *Client) GetPublic(
	address string,
	path string,
	result any,
) error {
	return c.do(
		http.MethodGet,
		address,
		path,
		nil,
		result,
		false,
	)
}

func (c *Client) Post(
	address string,
	path string,
	body any,
	result any,
) error {
	return c.do(
		http.MethodPost,
		address,
		path,
		body,
		result,
		true,
	)
}

func (c *Client) PostPublic(
	address string,
	path string,
	body any,
	result any,
) error {
	return c.do(
		http.MethodPost,
		address,
		path,
		body,
		result,
		false,
	)
}

func (c *Client) Put(
	address string,
	path string,
	body any,
	result any,
) error {
	return c.do(
		http.MethodPut,
		address,
		path,
		body,
		result,
		true,
	)
}

func (c *Client) Patch(
	address string,
	path string,
	body any,
	result any,
) error {
	return c.do(
		http.MethodPatch,
		address,
		path,
		body,
		result,
		true,
	)
}

func (c *Client) PatchPublic(
	address string,
	path string,
	body any,
	result any,
) error {
	return c.do(
		http.MethodPatch,
		address,
		path,
		body,
		result,
		false,
	)
}

func (c *Client) Delete(
	address string,
	path string,
	result any,
) error {
	return c.do(
		http.MethodDelete,
		address,
		path,
		nil,
		result,
		true,
	)
}

func (c *Client) DeletePublic(
	address string,
	path string,
	result any,
) error {
	return c.do(
		http.MethodDelete,
		address,
		path,
		nil,
		result,
		false,
	)
}

func (c *Client) do(
	method string,
	address string,
	path string,
	body any,
	result any,
	authenticated bool,
) error {
	var requestBody *bytes.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf(
				"encode request: %w",
				err,
			)
		}

		requestBody = bytes.NewReader(data)
	} else {
		requestBody = bytes.NewReader(nil)
	}

	url := strings.TrimRight(address, "/") +
		"/" +
		strings.TrimLeft(path, "/")

	req, err := http.NewRequest(
		method,
		url,
		requestBody,
	)
	if err != nil {
		return fmt.Errorf(
			"create request: %w",
			err,
		)
	}

	req.Header.Set(
		"Accept",
		"application/json",
	)

	if body != nil {
		req.Header.Set(
			"Content-Type",
			"application/json",
		)
	}

	if authenticated {
		if c.Token == "" {
			return errors.New(
				"authentication token is required",
			)
		}

		req.Header.Set(
			"Authorization",
			"Bearer "+c.Token,
		)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {
		return parseHTTPError(resp)
	}

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(
		resp.Body,
	).Decode(result); err != nil {
		return fmt.Errorf(
			"decode response: %w",
			err,
		)
	}

	return nil
}

func parseHTTPError(
	resp *http.Response,
) error {
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		message := strings.TrimSpace(
			string(body),
		)

		if message != "" {
			return errors.New(message)
		}
	}

	return fmt.Errorf(
		"daemon returned status %d",
		resp.StatusCode,
	)
}
