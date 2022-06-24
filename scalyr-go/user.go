package sdk

import (
	"context"
	"log"
)

type User struct {
	EmailAddress      string   `json:"emailAddress"`
	Permission        string   `json:"permission,omitempty"`
	AllowedSearch     string   `json:"allowedSearch,omitempty"`
	AllowedDashboards []string `json:"allowedDashboards,omitempty"`
	Groups            []string `json:"groups,omitempty"`
}

type InviteUserRequest struct {
	AuthParams
	User
}

type InviteUserResponse struct {
	APIResponse
}

func (scalyr *ScalyrConfig) InviteUser(ctx context.Context, email string, permission string, allowedSearch string, allowedDashboards []string, groups []string) error {
	request := &InviteUserRequest{}
	request.EmailAddress = email
	request.Permission = permission
	request.AllowedSearch = allowedSearch
	request.AllowedDashboards = allowedDashboards
	request.Groups = groups
	response := &InviteUserResponse{}
	err := NewRequest("POST", "/api/inviteUser", scalyr).withWriteConfig().jsonRequest(request).jsonResponse(ctx, response)
	log.Printf("%v", response)
	return err
}

type UpdateUserRequest = InviteUserRequest
type UpdateUserResponse = InviteUserResponse

func (scalyr *ScalyrConfig) UpdateUser(ctx context.Context, email string, permission string, allowedSearch string, allowedDashboards []string, groups []string) error {
	request := &UpdateUserRequest{}
	request.EmailAddress = email
	request.Permission = permission
	request.AllowedSearch = allowedSearch
	request.AllowedDashboards = allowedDashboards
	request.Groups = groups
	response := &UpdateUserResponse{}
	err := NewRequest("POST", "/api/editUserPermissions", scalyr).withWriteConfig().jsonRequest(request).jsonResponse(ctx, response)
	log.Printf("%v", response)
	return err

}

type RevokeUserRequest = InviteUserRequest
type RevokeUserResponse = InviteUserResponse

func (scalyr *ScalyrConfig) RevokeUser(ctx context.Context, email string) error {
	request := &RevokeUserRequest{}
	request.EmailAddress = email
	response := &RevokeUserResponse{}
	err := NewRequest("POST", "/api/revokeAccess", scalyr).withWriteConfig().jsonRequest(request).jsonResponse(ctx, response)
	if err != nil {
		return err
	}
	log.Printf("%v", response)
	return validateAPIResponse(&response.APIResponse, "Error Revoking User")
}
