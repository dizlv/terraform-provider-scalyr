package sdk

import ()

type Event struct {
	Thread string                 `json:"thread,omitempty"`
	Sev    int                    `json:"sev,omitempty"`
	Ts     string                 `json:"ts,omitempty"`
	Attrs  map[string]interface{} `json:"attrs"`
}

type Thread struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SessionInfo struct {
	ServerType string `json:"serverType,omitempty"`
	ServerID   string `json:"serverId,omitempty"`
}

type CreateEventsRequestParams struct {
	Session     string       `json:"session,omitempty"`
	SessionInfo *SessionInfo `json:"sessionInfo,omitempty"`
	Events      *[]Event     `json:"events"`
	Threads     *[]Thread    `json:"threads,omitempty"`
}

type CreateEventsRequest struct {
	AuthParams
	CreateEventsRequestParams
}

type CreateEventsResponse struct {
	APIResponse
}

func (scalyr *ScalyrConfig) SendEvent(event *Event, thread *Thread, session string, sessionInfo *SessionInfo) error {
	events := make([]Event, 0)
	threads := make([]Thread, 0)
	if event != nil {
		events = append(events, *event)
	}
	if thread != nil {
		threads = append(threads, *thread)
	}
	return scalyr.SendEvents(&events, &threads, session, sessionInfo)
}

func (scalyr *ScalyrConfig) SendEvents(events *[]Event, threads *[]Thread, session string, sessionInfo *SessionInfo) error {
	response := &CreateEventsResponse{}
	request := &CreateEventsRequest{}
	request.Session = session
	request.SessionInfo = sessionInfo
	request.Events = events
	request.Threads = threads
	err := NewRequest("POST", "/api/addEvents", scalyr).withWriteLog().jsonRequest(request).jsonResponse(response)
	return err
}
