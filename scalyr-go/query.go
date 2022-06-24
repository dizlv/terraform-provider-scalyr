package sdk

import (
	"context"
	"log"
)

type Query interface {
	Range(start string, end string) Query
	Fetch(ctx context.Context) ([]map[string]interface{}, error)
	HasMore() bool
	Size() float64
}

type BaseQuery struct {
	Query             string `json:"query"`
	ContinuationToken string `json:"continuationToken,omitempty"`
	Start             string `json:"startTime"`
	End               string `json:"endTime"`
	config            *ScalyrConfig
}

type LogQuery struct {
	BaseQuery
}

type LogQueryRequest struct {
	AuthParams
	LogQuery
}

/*
func (s *ScalyrConfig) NewLogQuery(query string) Query {
	lq := &LogQuery{}
	lq.Query=query
	lq.config = s
	return lq
}

func (q *LogQuery) Range(start string, end string) *Query {
	q.Start = start
	q.End = end
	return q
}
*/

type PowerQuery struct {
	BaseQuery
	last *PowerQueryResponse
}

type PowerQueryRequest struct {
	AuthParams
	PowerQuery
}

type PowerQueryResponse struct {
	APIResponse
	ContinuationToken string              `json:"continuationToken"`
	MatchingEvents    float64             `json:"matchingEvents"`
	OmittedEvents     float64             `json:"omittedEvents"`
	Columns           []map[string]string `json:"columns"`
	Values            [][]interface{}     `json:"values"`
}

func (scalyr *ScalyrConfig) NewPowerQuery(query string) Query {
	pq := &PowerQuery{}
	pq.Query = query
	pq.config = scalyr
	return pq
}

func (q *PowerQuery) Range(start string, end string) Query {
	q.Start = start
	q.End = end
	return q
}

func (q *PowerQuery) Size() float64 {
	return q.last.MatchingEvents
}

func (q *PowerQuery) Fetch(ctx context.Context) ([]map[string]interface{}, error) {

	response := &PowerQueryResponse{}
	request := &PowerQueryRequest{PowerQuery: *q}
	err := NewRequest("POST", "/api/powerQuery", q.config).withReadLog().jsonRequest(request).jsonResponse(ctx, response)
	if err != nil {
		return nil, err
	}
	q.last = response
	log.Printf("Response Matching Events: %v", response.MatchingEvents)
	log.Printf("Response Status: %v", response.Status)
	log.Printf("Response ContinuationToken: %v", response.ContinuationToken)
	log.Printf("Response Columns: %v", response.Columns)
	log.Printf("Response Values: %v", response.Values)
	flatResponse, _ := response.Flatten()
	log.Printf("Response: %v", flatResponse)
	return flatResponse, nil
}

func (q *PowerQuery) HasMore() bool {
	return false // not supported
	//return q.ContinuationToken != ""
}

func (q *PowerQueryResponse) Flatten() ([]map[string]interface{}, error) {
	log.Printf("Flatten")
	var results []map[string]interface{}
	for _, valueList := range q.Values {
		log.Printf("Record: %v", valueList)
		valueHash := make(map[string]interface{})
		for j, value := range valueList {
			valueHash[q.Columns[j]["name"]] = value
		}
		results = append(results, valueHash)
	}
	return results, nil
}
