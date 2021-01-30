package sdk

import (
	_ "encoding/json"
	"log"
	"strconv"
	"time"
)

type TokenTime time.Time

func (tt *TokenTime) UnmarshalJSON(b []byte) (err error) {
	r := string(b)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(tt) = time.Unix(q/1000, 0)
	return nil
}

func (tt TokenTime) Time() time.Time {
	return time.Time(tt).UTC()
}

func (tt TokenTime) String() string {
	return tt.Time().String()
}

type Token struct {
	Creator    string    `json:"creator"`
	Permission string    `json:"permission"`
	ID         string    `json:"id"`
	Label      string    `json:"label"`
	token      string    `json:"token"`
	CreateDate TokenTime `json:"createDate"`
}

type ListTokensRequest struct {
	AuthParams
}

type ListTokensResponseParams struct {
	Tokens []Token `json:tokens`
}

type ListTokensResponse struct {
	APIResponse
	ListTokensResponseParams
}

func (scalyr *ScalyrConfig) ListTokens() (*[]Token, error) {
	response := &ListTokensResponse{}
	err := NewRequest("POST", "/api/listTokens", scalyr).withReadConfig().withWriteConfig().jsonRequest(&ListTokensRequest{}).jsonResponse(response)
	log.Printf("%v", response)
	return &response.Tokens, err
}
