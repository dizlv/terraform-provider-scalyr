package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	httpClientTimeout = 15 * time.Second
)

func validateAPIResponse(response *APIResponse, message string) error {
	if response.Status != "success" {
		return fmt.Errorf("API Failure: %v - %v", message, response.Message)
	}
	return nil
}

func (ap *AuthParams) setToken(token string) {
	log.Printf("Set Token to: %v\n", token)
	ap.Token = token
}

type AuthParams struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

type APIRequest interface {
	setToken(token string)
}

type APIResponse struct {
	Message     string `json:"message"`
	Status      string `json:"status"`
	ResponseObj interface{}
}

type APITokenForDelegatingAccountRequest struct {
	DelegatingAccount string `json:"delegatingAccount"`
	TokenType         string `json:"logRead"`
}

type APITokenForDelegatingAccountResponse struct {
}

func (status *APIResponse) error(message string) error {
	return validateAPIResponse(status, message)
}

func getEnvWithDefault(envKey string, defaultValue string) string {
	v := os.Getenv(envKey)
	if v != "" {
		return v
	}
	return defaultValue
}

type Request struct {
	requestType   string
	request       interface{}
	uri           string
	apiKey        string
	supportedKeys []string
	config        *ScalyrConfig
	responseBody  []byte
	err           error
}

func NewRequest(requestType string, uri string, config *ScalyrConfig) *Request {
	return &Request{requestType: requestType, uri: uri, config: config}
}

func (r *Request) withWriteLog() *Request {
	if r.apiKey != "" {
		return r
	}
	if r.config.hasTeam() {

	}
	if r.config.Tokens.WriteLog != "" {
		r.apiKey = r.config.Tokens.WriteLog
	} else {
		r.supportedKeys = append(r.supportedKeys, "WriteLog")
	}
	return r
}

func (r *Request) withReadLog() *Request {
	if r.apiKey != "" {
		return r
	}
	if r.config.Tokens.ReadLog != "" {
		r.apiKey = r.config.Tokens.ReadLog
	} else {
		r.supportedKeys = append(r.supportedKeys, "ReadLog")
	}
	return r
}

func (r *Request) withReadConfig() *Request {
	if r.apiKey != "" {
		return r
	}
	if r.config.Tokens.ReadConfig != "" {
		r.apiKey = r.config.Tokens.ReadConfig
	} else {
		r.supportedKeys = append(r.supportedKeys, "ReadConfig")
	}
	return r
}

func (r *Request) withWriteConfig() *Request {
	if r.apiKey != "" {
		return r
	}
	if r.config.Tokens.WriteConfig != "" {
		r.apiKey = r.config.Tokens.WriteConfig
	} else {
		r.supportedKeys = append(r.supportedKeys, "WriteConfig")
	}
	return r
}

func (r *Request) jsonRequest(request APIRequest) *Request {
	r.request = request
	return r
}

func (r *Request) emptyRequest() *Request {
	r.request = APIRequest(&AuthParams{})
	return r
}

func prepareRequestBody(data interface{}) ([]byte, error) {
	switch data := data.(type) {
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		return b, nil
	}
}

func (r *Request) jsonResponse(response interface{}) error {
	if r.err != nil {
		return r.err
	}

	if r.request == nil {
		r.emptyRequest()
	}

	if r.apiKey == "" && len(r.supportedKeys) > 0 {
		return fmt.Errorf("No API Key Found - Supported Tokens for %v are %v", r.uri, r.supportedKeys)
	} else {
		r.request.(APIRequest).setToken(r.apiKey)
	}

	// Validate endpoint
	parsedUrl, err := url.Parse(r.config.Endpoint)
	if err != nil {
		return err
	}

	parsedUrl.Path = path.Join(parsedUrl.Path, r.uri)
	body, err := prepareRequestBody(r.request)

	request, err := http.NewRequest(r.requestType, parsedUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		return err
	}

	log.Printf("Outgoing Request: %s", dumpRequest)

	httpClient := http.Client{
		Timeout: httpClientTimeout,
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	dumpResponse, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	log.Printf("Incoming Response: %s", dumpResponse)

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.Unmarshal(responseBody, &response); err != nil {
			return fmt.Errorf("could not decode response JSON, %s: %v", string(responseBody), err)
		}

		return nil

	default:
		return fmt.Errorf("")
	}
}
