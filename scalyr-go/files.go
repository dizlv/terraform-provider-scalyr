package sdk

import (
	"fmt"
)

type GetFileRequest struct {
	AuthParams
	Path            string `json:"path"`
	ExpectedVersion int64  `json:"expectedVersion,omitempty"`
}

type GetFileResponse struct {
	APIResponse
	Path       string  `json:"path"`
	Version    int64   `json:"version"`
	CreateDate APITime `json:"createDate"`
	ModDate    APITime `json:"modDate"`
	Content    string  `json:"content"`
}

func (scalyr *ScalyrConfig) GetFile(path string) (*GetFileResponse, error) {
	getFileResponse := &GetFileResponse{}
	getFileRequest := &GetFileRequest{Path: path}

	err := NewRequest("POST", "/api/getFile", scalyr).withWriteConfig().withReadConfig().jsonRequest(getFileRequest).jsonResponse(getFileResponse)
	return getFileResponse, err
}

type PutFileRequest struct {
	AuthParams
	Path            string `json:"path"`
	Content         string `json:"content"`
	ExpectedVersion int64  `json:"expectedVersion,omitempty"`
}

type PutFileResponse struct {
	APIResponse
	Path string `json:"path"`
}

func (scalyr *ScalyrConfig) PutFile(path string, content string) (*PutFileResponse, error) {
	response := &PutFileResponse{}
	request := &PutFileRequest{Path: path, Content: content}

	err := NewRequest("POST", "/api/putFile", scalyr).withWriteConfig().jsonRequest(request).jsonResponse(response)
	return response, err
}

type DeleteFileRequest struct {
	AuthParams
	Path            string `json:"path"`
	DeleteFile      bool   `json:"deleteFile,omitempty"`
	ExpectedVersion int64  `json:"expectedVersion,omitempty"`
}

type DeleteFileResponse struct {
	APIResponse
}

func (scalyr *ScalyrConfig) DeleteFile(path string) error {
	response := &DeleteFileResponse{}
	request := &DeleteFileRequest{Path: path, DeleteFile: true}

	err := NewRequest("POST", "/api/putFile", scalyr).withWriteConfig().jsonRequest(request).jsonResponse(response)
	if response.Status != "success" {
		return fmt.Errorf("Error Deleting %v - %v", path, response.Status)
	}
	return err
}
