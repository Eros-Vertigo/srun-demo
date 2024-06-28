package slzx

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"github.com/dot-xiaoyuan/srun-demo/floger"
	"net/http"
	"net/url"
)

// 发送请求
func doRequest(uri, method string, data interface{}) (*http.Response, error) {
	client := &http.Client{
		// 忽略证书验证
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	var httpReq *http.Request
	baseURL := uri
	if method == http.MethodPost {
		httpReq, _ = http.NewRequest(method, baseURL, bytes.NewBuffer(data.([]byte)))
	} else {
		httpReq, _ = http.NewRequest(method, baseURL, nil)
		params := data.(url.Values)
		httpReq.URL.RawQuery = params.Encode()
		floger.Debug5("Request URL:", httpReq.URL.String())
	}

	// Set the content type to application/type
	httpReq.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	return resp, nil
}

func GenerateChannelSerialNumber() string {
	// Generate 16 random bytes
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	// Encode bytes to hexadecimal string
	randomString := hex.EncodeToString(randomBytes)
	return randomString
}
