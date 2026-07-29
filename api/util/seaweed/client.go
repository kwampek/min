package seaweed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type SeaweedClient struct {
	MasterURL string
	VolumeURL string

	Client *http.Client
}

func MakeSeaweedClient(masterURL, volumeURL string) SeaweedClient {
	return SeaweedClient{
		MasterURL: masterURL,
		VolumeURL: volumeURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type assignResponse struct {
	Fid string `json:"fid"`
}

type uploadResponse struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func (c *SeaweedClient) Upload(r io.Reader, filename string) (string, int64, error) {
	resp, err := c.Client.Get(c.MasterURL + "/dir/assign")
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("assign failed: %s", string(body))
	}

	var assign assignResponse
	if err := json.NewDecoder(resp.Body).Decode(&assign); err != nil {
		return "", 0, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if filename == "" {
		filename = "file"
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", 0, err
	}

	size, err := io.Copy(part, r)
	if err != nil {
		return "", 0, err
	}

	if err := writer.Close(); err != nil {
		return "", 0, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", c.VolumeURL, assign.Fid),
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

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("upload failed: %s", string(b))
	}

	return assign.Fid, size, nil
}

func (c *SeaweedClient) Open(fid string) (io.ReadCloser, error) {
	// using only one volume
	resp, err := c.Client.Get(fmt.Sprintf("%s/%s", c.VolumeURL, fid))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("file not found")
	}

	return resp.Body, nil
}

func (c *SeaweedClient) Delete(fid string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%s", c.VolumeURL, fid),
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
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(b))
	}

	return nil
}
