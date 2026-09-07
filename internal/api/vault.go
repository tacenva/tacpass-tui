package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func (c *Client) ListVaults(
	address string,
	result any,
) error {
	debugFile, err := os.OpenFile(
		"debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err == nil {
		defer debugFile.Close()

		fmt.Fprintf(
			debugFile,
			"ListVaults: GET %s/vault\nToken exists: %t\n",
			strings.TrimRight(address, "/"),
			c.Token != "",
		)
	}

	return c.Get(
		address,
		"/vault",
		result,
	)
}

func (c *Client) CreateVault(
	address string,
	body any,
	result any,
) error {
	return c.Post(
		address,
		"/vault",
		body,
		result,
	)
}

func (c *Client) ListRecords(
	address string,
	vaultID string,
	result any,
) error {
	path := fmt.Sprintf(
		"/vault/%s/record",
		vaultID,
	)

	return c.Get(
		address,
		path,
		result,
	)
}

func (c *Client) GetRecord(
	address string,
	vaultID string,
	recordID string,
	result any,
) error {
	path := fmt.Sprintf(
		"/vault/%s/record/%s",
		vaultID,
		recordID,
	)

	return c.Get(
		address,
		path,
		result,
	)
}

// CreateRecordRaw mengirim encrypted record sebagai raw bytes.
//
// Tidak menggunakan Post() karena Post() akan melakukan json.Marshal()
// terhadap []byte dan menghasilkan base64 JSON.
func (c *Client) CreateRecordRaw(
	address string,
	vaultID string,
	data []byte,
) (string, error) {
	path := fmt.Sprintf(
		"/vault/%s/record",
		vaultID,
	)

	body, err := c.doRaw(
		http.MethodPost,
		address,
		path,
		data,
	)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// UpdateRecordRaw mengirim encrypted record sebagai raw bytes.
//
// Tidak menggunakan Put() karena Put() akan melakukan json.Marshal()
// terhadap []byte.
func (c *Client) UpdateRecordRaw(
	address string,
	vaultID string,
	recordID string,
	data []byte,
) error {
	path := fmt.Sprintf(
		"/vault/%s/record/%s",
		vaultID,
		recordID,
	)

	_, err := c.doRaw(
		http.MethodPut,
		address,
		path,
		data,
	)

	return err
}

func (c *Client) DeleteRecord(
	address string,
	vaultID string,
	recordID string,
) error {
	path := fmt.Sprintf(
		"/vault/%s/record/%s",
		vaultID,
		recordID,
	)

	return c.Delete(
		address,
		path,
		nil,
	)
}

func (c *Client) doRaw(
	method string,
	address string,
	path string,
	data []byte,
) ([]byte, error) {
	url := strings.TrimRight(address, "/") +
		"/" +
		strings.TrimLeft(path, "/")

	req, err := http.NewRequest(
		method,
		url,
		strings.NewReader(string(data)),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create request: %w",
			err,
		)
	}

	req.Header.Set(
		"Accept",
		"application/octet-stream",
	)

	req.Header.Set(
		"Content-Type",
		"application/octet-stream",
	)

	if c.Token == "" {
		return nil, fmt.Errorf(
			"authentication token is required",
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+c.Token,
	)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseHTTPError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"read response: %w",
			err,
		)
	}

	return body, nil
}
