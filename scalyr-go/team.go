package sdk

import "context"

type Team struct {
	Token        string `json:"token"`
	EmailAddress string `json:"emailAddress"`
}

type CreateTeamRequestParams struct {
	EmailAddress string `json:"emailAddress"`
}

type CreateTeamRequest struct {
	AuthParams
	CreateTeamRequestParams
}

type CreateTeamResponseParams struct {
	Token string `json:"token"`
}

type CreateTeamResponse struct {
	APIResponse
	CreateTeamResponseParams
}

func (scalyr *ScalyrConfig) CreateTeam(ctx context.Context, email string) (*CreateTeamResponse, error) {
	response := &CreateTeamResponse{}
	err := NewRequest("POST", "/api/listTeamAccounts", scalyr).withReadConfig().withWriteConfig().jsonRequest(&ListTeamsRequest{}).jsonResponse(ctx, response)
	return response, err
}

type ListTeamsRequest struct {
	AuthParams
}

type ListTeamsResponse struct {
	APIResponse
	Teams []string `json:"teams"`
}

func (scalyr *ScalyrConfig) ListTeams(ctx context.Context) (*[]string, error) {
	teamsResponse := &ListTeamsResponse{}
	err := NewRequest("POST", "/api/listTeamAccounts", scalyr).withReadConfig().withWriteConfig().jsonRequest(&ListTeamsRequest{}).jsonResponse(ctx, teamsResponse)
	return &teamsResponse.Teams, err
}
