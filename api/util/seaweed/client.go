package seaweed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type SeaweedClient struct {
	MasterURL string

	Client *http.Client
}

func New(masterURL string) *SeaweedClient {
	return &SeaweedClient{
		MasterURL: masterURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type assignResponse struct {
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicURL string `json:"publicUrl"`
}

type uploadResponse struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Error string `json:"error"`
}

func (c *SeaweedClient) Upload(
	ctx context.Context,
	filename string,
	r io.Reader,
) (fid string, size int64, err error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.MasterURL+"/dir/assign",
		nil,
	)
	if err != nil {
		return "", 0, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var assign assignResponse
	if err := json.NewDecoder(resp.Body).Decode(&assign); err != nil {
		return "", 0, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", 0, err
	}

	n, err := io.Copy(part, r)
	if err != nil {
		return "", 0, err
	}

	if err := writer.Close(); err != nil {
		return "", 0, err
	}

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://"+assign.Url+"/"+assign.Fid,
		body,
	)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err = c.Client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var upload uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&upload); err != nil {
		return "", 0, err
	}

	if upload.Error != "" {
		return "", 0, fmt.Errorf(upload.Error)
	}

	return assign.Fid, n, nil
}

func (c *SeaweedClient) Open(
	ctx context.Context,
	fid string,
) (io.ReadCloser, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.VolumeURL+"/"+fid,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("seaweed returned %s", resp.Status)
	}

	return resp.Body, nil
}

func (c *SeaweedClient) Delete(
	ctx context.Context,
	fid string,
) error {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		c.VolumeURL+"/"+fid,
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted &&
		resp.StatusCode != http.StatusOK {

		return fmt.Errorf("seaweed returned %s", resp.Status)
	}

	return nil
}
