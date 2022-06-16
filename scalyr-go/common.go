package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
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

	url := fmt.Sprintf("%v%v", r.config.Endpoint, r.uri)
	buffer, err := json.Marshal(r.request)
	if err != nil {
		r.err = err
		return err
	}
	log.Printf("Call %v %v with %v\n", r.requestType, url, string(buffer))
	req, err := http.NewRequest(r.requestType, url, bytes.NewReader(buffer))
	if err != nil {
		log.Printf("Failure in creating NewRequest")
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	httpClient := http.Client{}
	request, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}

	log.Printf("Request: %s", request)

	resp, err := httpClient.Do(req)
	if err != nil {
		r.err = err
		log.Printf("Error %v", r.err)
		return err
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("RESPONSE:\n%s", string(respDump))

		return fmt.Errorf("wrong status code %d", resp.StatusCode)
	}

	r.responseBody, r.err = ioutil.ReadAll(resp.Body)
	log.Printf("Response")
	for _, chunk := range Chunk(string(r.responseBody), 200) {
		log.Printf("%v", chunk)
	}
	return json.Unmarshal(r.responseBody, &response)
}
